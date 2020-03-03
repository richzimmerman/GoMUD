package client

import (
	"account"
	"bufio"
	"fmt"
	"mobs"
	"net"
	"output"
	"sync"
	"telnet"
	"utils"
)

const (
	statePrompt                       = 0
	stateLogin                        = 1
	stateLoggedIn                     = 2
	stateLogout                       = 3
	stateLoggedOut                    = 4
	stateAccount                      = 5
	statePassword                     = 6
	stateDisconnected                 = 7
	stateAccountMenu                  = 8
	stateAccountNewCharacter          = 9
	stateAccountListCharacters        = 10
	stateAccountDeleteCharacter       = 11
	stateAccountChangePassword        = 12
	stateAccountQuit                  = 13
	stateCreateAccount                = 14
	stateCreateAccountName            = 15
	stateCreateAccountNameConfirm     = 16
	stateCreateAccountPassword        = 17
	stateCreateAccountPasswordConfirm = 18
	stateCreateAccountEmail           = 19
	stateVerify                       = 20
	stateInGame                       = 21

	accountPrompt  = "<Y>Account:</Y> "
	passwordPrompt = "<Y>Password:</Y> "

	invalidAccount  = "Account does not exist. Create new account? (Y/n)\n"
	invalidPassword = "Invalid password, please try logging in again.\n"

	accountMenu = `
<Y>Account Menu</Y>
<BW>------------</BW>
<Y>N)</Y> New Character
<Y>L)</Y> List Characters
<Y>D)</Y> Delete Character
<Y>C)</Y> Change Password
<Y>Q)</Y> Quit
<BW>------------</BW>
<Y>Enter a menu item or character name to log in:</Y>
`
)

type Client struct {
	loggedIn    bool
	state       int8
	Connection  net.Conn
	Telnet      *telnet.Telnet
	Name        string
	Player      *mobs.Player
	Account     *account.Account
	OutputSteam chan string
	outMutex    sync.Mutex
	In          *bufio.Reader
}

func NewClient(c net.Conn) *Client {
	in := bufio.NewReader(c)
	client := &Client{
		Connection:  c,
		Telnet:      telnet.NewTelnet(c, in),
		In:          in,
		OutputSteam: make(chan string),
	}
	return client
}

func (c *Client) outListener() {
	/*
		This is the main output listener go routine that will right output to the client w/ mutex to help try and avoid
		race conditions
	*/
	go func() {
		defer close(c.OutputSteam)

		for {
			s, ok := <-c.OutputSteam
			if !ok {
				fmt.Println("unable to read out channel")
			}
			bytes, err := output.ANSIFormatter(s)
			if err != nil {
				fmt.Printf("failed to parse output string: %v\n", err)
				continue
			}
			c.outMutex.Lock()
			_, err = c.Connection.Write(bytes)
			c.outMutex.Unlock()
			if err != nil {
				fmt.Printf("unable to write to client connection: %v\n", err)
				return
			}
		}

	}()
}

func (c *Client) GameLoop() error {
	/*
		The main game loop for the client, which will go through the login, in game and logout processes
	*/
	defer c.Connection.Close()

	c.outListener()

	c.state = stateLogin
	for {
		switch c.state {
		case statePrompt:
			break
		case stateLogin:
			ok, err := c.logIn()
			if err != nil {
				return utils.Error(fmt.Errorf("failed login: %v", err))
			}
			if ok {
				c.state = stateLoggedIn
			} else {
				c.state = stateDisconnected
			}
			break
		case stateLoggedIn:
			err := c.accountMenu()
			if err != nil {
				return utils.Error(err)
			}
			break
		case stateInGame:
			c.inGameLoop()
			break
		case stateLogout:
			// TODO: start logout
			//c.logout()
			break
		case stateLoggedOut:
			// TODO: go back to login
			c.state = stateLogin
			break
		case stateDisconnected:
			return nil
		}
	}
}

func (c *Client) inGameLoop() {
	/*
		The main in-game loop, starts after login
	*/
	c.OutputSteam <- "<R>You've logged into the game! Huzzah!</R>"
	for c.state == stateInGame {
		_, err := c.Telnet.Read()
		if err != nil {
			// If we're unable to read from connection, the socket is likely broken from the client side.
			c.state = stateDisconnected
		}
	}
}
