package telnet

import (
	"bufio"
	"bytes"
	"fmt"
	"input"
	"lib/commands"
	"lib/sessions"
	"logger"
	"net"
	"strings"
	"time"
	"utils"
)

const (
	MaxInputSize = 128
	BufferSize   = 2048

	iac  = byte(255)
	will = byte(251)
	wont = byte(252)
	do   = byte(253)
	dont = byte(254)
	sb   = byte(250)
	se   = byte(240)

	optBinary = byte(0)
	optEcho   = byte(1)
	optGmcp   = byte(201)
	optClient = byte(222)

	stateData         = 0
	stateIac          = 1
	stateIacSb        = 2
	stateIacWill      = 3
	stateIacDo        = 4
	stateIacWont      = 5
	stateIacDont      = 6
	stateIacSbIac     = 7 // This one is likely not used
	stateIacSbData    = 8
	stateIacSbDataIac = 9
)

var log = logger.NewLogger()

type Telnet struct {
	Connection     net.Conn
	InputSteam     *bufio.Reader
	Data           chan []byte
	subnegOffset   int
	subnegotiation []byte
	negState       int
	commandQueue   chan *input.Input
}

func NewTelnet(c net.Conn, i *bufio.Reader) *Telnet {
	t := &Telnet{
		Connection:   c,
		InputSteam:   i,
		Data:         make(chan []byte),
		subnegOffset: 0,
		commandQueue: make(chan *input.Input),
	}
	go t.StartParsing()
	return t
}

// Very similar to the Read function, but this returns the input string for a prompt rather than
// continually feeding input into the input queue, but will still handle subnegotiations while waiting for input.
func (t *Telnet) Prompt() (string, error) {
	i := 0
	inputBuffer := make([]byte, MaxInputSize)
	for i == 0 {
		tempIn := make([]byte, BufferSize)
		length, err := t.InputSteam.Read(tempIn)
		if err != nil {
			return "", utils.Error(err)
		}
		i, err = t.negotiate(inputBuffer, tempIn)
		if length < 0 {
			i = length
			break
		}
	}
	// Remove trailing null bytes
	inputBuffer = bytes.Trim(inputBuffer, "\x00")
	input := string(inputBuffer)
	// Remove whitespace to clean up input
	input = strings.TrimSpace(input)
	return input, nil
}

func (t *Telnet) Read() (int, error) {
	/*
		This function shouldn't be used until after a user has logged into a character. So the *mobs.Player
		should not be nil, and user input is generic enough to be passed along to the parser
	*/
	i := 0
	inputBuffer := make([]byte, MaxInputSize)
	for i == 0 {
		tempIn := make([]byte, BufferSize)
		length, err := t.InputSteam.Read(tempIn)
		if err != nil {
			return -1, utils.Error(err)
		}
		i, err = t.negotiate(inputBuffer, tempIn)
		if length < 0 {
			i = length
			break
		}
	}
	// Remove trailing null bytes
	inputBuffer = bytes.Trim(inputBuffer, "\x00")
	inputString := string(inputBuffer)
	// Trim whitespace (new lines)
	s, _, err := sessions.GetSessionByIpAddress(t.Connection.RemoteAddr().String())
	if err != nil {
		s.Client().Out("You cannot do that because you shouldn't exist!!")
	} else {
		parsableInput := input.NewInput(s, strings.TrimSpace(inputString))
		t.commandQueue <- parsableInput
	}
	return i, nil
}

func (t *Telnet) StartParsing() {
	t.commandQueue = make(chan *input.Input)
	defer close(t.commandQueue)

	for {
		select {
		case input := <-t.commandQueue:
			s, _, e := sessions.GetSessionByIpAddress(t.Connection.RemoteAddr().String())
			if e != nil {
				log.Err("%v", e)
			}
			s.InputReceived(time.Now().Unix())
			// It's probably important to NOT parse inputs concurrently to prevent possible race conditions
			err := t.parseInput(input)
			if err != nil {
				s.Client().Out("Command not found.")
			}
			break
		}
	}
}

func (t *Telnet) parseInput(i *input.Input) error {
	// TODO: chat parsing before general command parsing (PROBABLY NOT THOUGH)
	split := strings.Split(i.Input(), " ")
	cmd := split[0]
	command, err := commands.GetCommand(cmd)
	if err != nil {
		return utils.Error(err)
	}
	command.Execute(i.Session(), split)
	return nil
}

func (t *Telnet) respondToNegotiation(response byte, opt byte) (err error) {
	optResponse := []byte{iac, response, opt}
	_, err = t.Connection.Write(optResponse)
	if err != nil {
		fmt.Println(err)
	}
	return nil
}

func (t *Telnet) negotiate(outBuffer []byte, inData []byte) (n int, err error) {
	reply := byte(0)
	offset := 0
	exit := false
	for i := 0; i < len(inData); i++ {
		b := inData[i]
		if b == 0 {
			exit = true
		}
		if exit {
			break
		}
		if i >= MaxInputSize {
			// To prevent a buffer overflow, silently eat the input and move on
			return 0, nil
		}
		switch t.negState {
		case stateData:
			if b == iac {
				t.negState = stateIac
			} else {
				outBuffer[offset] = b
				offset++
			}
			break
		case stateIac:
			switch b {
			case iac:
				t.negState = stateData
				outBuffer[i] = iac
				break
			case will:
				t.negState = stateIacWill
				break
			case wont:
				t.negState = stateIacWont
				break
			case dont:
				t.negState = stateIacDont
				break
			case do:
				t.negState = stateIacDo
				break
			case sb:
				t.negState = stateIacSb
				break
			case se:
				exit = true
				break
			default:
				t.negState = stateData
				break
			}
			break
		case stateIacWill:
			switch b {
			case optEcho:
				reply = dont
				break
			case optBinary:
			case optGmcp:
			case optClient:
				reply = do
				break
			default:
				reply = dont
				break
			}
			err = t.respondToNegotiation(reply, b)
			if err != nil {
				fmt.Println(err)
			}
			t.negState = stateData
			break
		case stateIacWont:
			switch b {
			case optClient:
				reply = do
				break
			default:
				reply = dont
				break
			}
			err = t.respondToNegotiation(reply, b)
			if err != nil {
				fmt.Println(err)
			}
			t.negState = stateData
			break
		case stateIacDo:
			reply = wont
			err = t.respondToNegotiation(reply, b)
			if err != nil {
				fmt.Println(err)
			}
			t.negState = stateData
			break
		case stateIacDont:
			reply = wont
			err = t.respondToNegotiation(reply, b)
			if err != nil {
				fmt.Println(err)
			}
			t.negState = stateData
			break
		case stateIacSb:
			t.negState = stateIacSbData
			break
		case stateIacSbData:
			switch b {
			case iac:
				t.negState = stateIacSbDataIac
				break
			default:
				t.subnegotiation[t.subnegOffset] = b
				t.subnegOffset++
				break
			}
			break
		case stateIacSbDataIac:
			switch b {
			case se:
				// Handle the GMCP messsage and empty the subnegotiation buffer
				t.Data <- bytes.Trim(t.subnegotiation, "\x00")
				t.subnegotiation = make([]byte, 2048)
				t.subnegOffset = 0
				t.negState = stateData
				break
			default:
				t.negState = stateData
				break
			}
			break

		}
	}
	if offset == 0 {
		return -1, err
	}
	return offset, err
}
