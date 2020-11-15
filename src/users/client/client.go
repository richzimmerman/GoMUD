package client

import (
	"bufio"
	"fmt"
	. "lib/players"
	"lib/sessions"
	. "lib/world"
	"logger"
	"message"
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

	accountAlreadyLoggedIn = "That account is already logged in."
	invalidAccount         = "Account does not exist. Create new account? (Y/n)\n"
	invalidPassword        = "Invalid password, please try logging in again.\n"

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

var log = logger.NewLogger()

type Client struct {
	loggedIn     bool
	state        int8
	Connection   net.Conn
	Telnet       *telnet.Telnet
	Name         string
	OutputStream chan string
	outMutex     sync.Mutex
	In           *bufio.Reader
	AccountInfo  *loggedInfo
}

type loggedInfo struct {
	Account string
	Player  string
}

// NewClient returns a new client struct for an associated connection
func NewClient(c net.Conn) *Client {
	in := bufio.NewReader(c)
	client := &Client{
		Connection:   c,
		Telnet:       telnet.NewTelnet(c, in),
		In:           in,
		OutputStream: make(chan string),
		AccountInfo:  &loggedInfo{Account: "", Player: ""},
	}
	return client
}

// GetRemoteAddress returns the remote IP address (including the port) for a direct reference to a client
// Includes the port so multiple connections from the same IP address are not blocked
func (c *Client) GetRemoteAddress() string {
	return c.Connection.RemoteAddr().String()
}

// AssociatedAccount returns the name of the account associated with this logged in session
func (c *Client) AssociatedAccount() (string, error) {
	var err error
	if c.AccountInfo.Account == "" && c.loggedIn {
		err = fmt.Errorf("client is not logged into an account")
	}
	return c.AccountInfo.Account, err
}

// SetAssociatedAccount sets account name associated with this client connection
func (c *Client) SetAssociatedAccount(name string) {
	c.AccountInfo.Account = name
}

// AssociatedPlayer returns the name of the currently logged in player for this session
func (c *Client) AssociatedPlayer() string {
	return c.AccountInfo.Player
}

// SetAssociatedPlayer sets the name of the player struct associated with this connection
func (c *Client) SetAssociatedPlayer(name string) {
	log.Info("setting player %s", name)
	c.AccountInfo.Player = name
}

func (c *Client) Out(msg string) {
	c.OutputStream <- msg
}

func (c *Client) Logout() {
	// TODO: handle this in sessions... but call the sessions kill handler
	// This is nil
	p, err := GetPlayer(c.AssociatedPlayer())
	if err != nil {
		log.Err("killing client without associated player: %v", err)
		return
	}
	//
	r, err := GetRoom(p.GetLocation())
	if err != nil {
		log.Err("unable to logout %s from room (%s): %v", p.GetName(), p.GetLocation(), err)
		return
	}
	r.RemovePlayer(p.GetName())
	RemovePlayer(p.GetName())
	c.Connection.Close()
}

func (c *Client) outListener() {
	/*
		This is the main output listener go routine that will right output to the client w/ mutex to help try and avoid
		race conditions
	*/
	go func() {
		defer close(c.OutputStream)

		for {
			s, ok := <-c.OutputStream
			if !ok {
				log.Err("unable to read out channel: %s", c.Connection.LocalAddr())
				return
			}
			bytes, err := output.ANSIFormatter(s)
			if err != nil {
				log.Err("failed to parse output string: %v\n", err)
				continue
			}
			c.outMutex.Lock()
			_, err = c.Connection.Write(bytes)
			c.outMutex.Unlock()
			if err != nil {
				log.Err("unable to write to client connection: %v\n", err)
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

	var err error
	var accountName string
	c.state = stateLogin
	for {
		switch c.state {
		case statePrompt:
			break
		case stateLogin:
			accountName, err = c.logIn()
			if err != nil {
				return utils.Error(fmt.Errorf("failed login: %v", err))
			}
			if accountName != "" {
				c.state = stateLoggedIn
			} else {
				c.state = stateDisconnected
			}
			break
		case stateLoggedIn:
			err := c.accountMenu(accountName)
			if err != nil {
				return utils.Error(err)
			}
			break
		case stateInGame:
			log.Info("starting game loop with player: %s", c.AssociatedPlayer())
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
	player, err := CheckPlayer(c.AssociatedPlayer(), c)
	if err != nil {
		log.Err("%v", err)
		log.Err("in game loop started for %s without an associated player", c.Connection.LocalAddr().String())
		c.Connection.Close()
	}
	if player.GetLocation() == "" {
		_, err := GetRoom("0")
		if err != nil {
			log.Err("the world is empty and does not have a 'limbo' room")
		}
		player.SetLocation("0")
	}
	room, err := GetRoom(player.GetLocation())
	if err != nil {
		log.Err("room (%s) does not seem to exist", player.GetLocation())
		c.Connection.Close()
	}
	sessions.CreateSession(c, player)
	// TODO: Player does not exist in room, might need to add player to room rather than move
	MovePlayer(player, player.GetLocation())
	u := message.NewUnformattedMessage("You appear out of thin air!", "", "<Y><A.NAME></Y> appears out of thin air!")
	msg := message.NewMessage(player, nil, u)
	room.Send(msg)
	c.Out(room.Look(player))
	for c.state == stateInGame {
		i, err := c.Telnet.Read()
		if err != nil || i < 0 {
			// If we're unable to read from connection, the socket is likely broken from the client side.
			c.state = stateDisconnected
		}
	}
}
