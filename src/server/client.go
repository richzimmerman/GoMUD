package server

import (
	"account"
	"bufio"
	"fmt"
	"mobs"
	"net"
	"output"
	"strings"
	"sync"
	"telnet"

	"db"
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
	stateVerify                       = 19

	accountPrompt  = "<G>Account:</G> "
	passwordPrompt = "<Bl>Password:</Bl> "

	invalidAccount  = "Account does not exist.\n"
	invalidPassword = "Invalid password, please try logging in again.\n"

	accountMenu = `
<Y>Account Menu</Y>
<BW>------------</BW>
<Y>N)</Y> <W>New Character</W>
<Y>L)</Y> <W>List Characters</W>
<Y>D)</Y> <W>Delete Character</W>
<Y>C)</Y> <W>Change Password</W>
<Y>Q)</Y> <W>Quit</W>
<BW>------------</BW>
<Y>Enter a menu item or character name to log in:</Y>
`
)

type Client struct {
	loggedIn    bool
	state       int8
	Connection  net.Conn
	Address     string
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
		Address:     c.RemoteAddr().String(),
		Telnet:      telnet.NewTelnet(c, in),
		In:          in,
		OutputSteam: make(chan string),
	}
	return client
}

func (c *Client) outListener() {
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

func (c *Client) prompt() (string, error) {
	/* Generic Prompt method to get a single input for menu/login related prompts */
	for {
		input, err := c.In.ReadString('\n')
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(input), nil
	}
}

func (c *Client) gameLoop() error {
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
				return fmt.Errorf("unable to login: %v \n", err)
			}
			if ok {
				c.state = stateLoggedIn
			}
			break
		case stateLoggedIn:
			err := c.accountMenu()
			if err != nil {
				return err
			}
			//c.mainLoop()
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

func (c *Client) mainLoop() {
	for {
		_, err := c.Telnet.Read()
		if err != nil {
			//fmt.Printf("error reading from connection, killing it: %v \n", err)
			c.state = stateDisconnected
			break
		}

	}
}

func (c *Client) logIn() (bool, error) {
	var accountName string
	var err error
	loginState := stateAccount
	c.state = stateAccount
	for {
		switch c.state {
		case statePrompt:
			switch loginState {
			case stateAccount:
				accountName, err = c.prompt()
				if err != nil {
					return false, err
				}
				ok, err := db.DatabaseConnection.AccountExists(accountName)
				if err != nil {
					return false, err
				}
				if ok {
					loginState = statePassword
					c.state = statePassword
				} else {
					c.OutputSteam <- invalidAccount
					c.state = stateAccount
				}
				break
			case statePassword:
				pw, err := c.prompt()
				if err != nil {
					return false, err
				}
				ok, err := db.DatabaseConnection.VerifyPassword(accountName, pw)
				if err != nil {
					return false, err
				}
				if ok {
					c.Account, err = account.LoadAccount(accountName)
					if err != nil {
						return false, err
					}
					return true, nil
				} else {
					c.OutputSteam <- invalidPassword
					loginState = stateAccount
					c.state = stateAccount
				}
				break
			}
			break
		case stateAccount:
			c.OutputSteam <- accountPrompt
			c.state = statePrompt
			break
		case statePassword:
			c.OutputSteam <- passwordPrompt
			c.state = statePrompt
			break
		}
	}
}

//func (c *Client) createAccountPrompt() (bool, error) {
//	var accountName string
//	var password string
//	var err error
//
//	c.state = stateCreateAccount
//	var createState = stateCreateAccount
//	for {
//		switch c.state {
//		case statePrompt:
//			switch createState {
//			case stateCreateAccount:
//
//			}
//		}
//	}
//}

func (c *Client) createAccount(accountName string, password string, email string) error {
	lastip := strings.SplitN(c.Address, ":", 1)
	fmt.Printf("last ip: %s \n", lastip)
	a, err := account.NewAccount(accountName, password, lastip[0], email)
	if err != nil {
		return fmt.Errorf("unable to create account %s: %v", accountName, email)
	}
	c.Account = a
	return nil
}

func (c *Client) accountMenu() error {
	var input string
	var err error
	c.state = stateAccountMenu
	for {
		switch c.state {
		case statePrompt:
			input, err = c.prompt()
			if err != nil {
				return fmt.Errorf("unable to login via account menu: %v", err)
			}
			switch strings.ToLower(input) {
			case "n":
				c.state = stateAccountNewCharacter
				break
			case "l":
				c.state = stateAccountListCharacters
				break
			case "d":
				c.state = stateAccountDeleteCharacter
				break
			case "c":
				c.state = stateAccountChangePassword
				break
			case "q":
				c.state = stateAccountQuit
				break
			default:
				c.OutputSteam <- fmt.Sprintf("<Y>You've enterred:</Y> %s", input)
				// TODO: accept character name input
				// ok, err := c.Account.VerifyCharacter(input)
				//if err != nil {
				//	return fmt.Errorf("unable to verify character: %s for account: %s", input, c.Account.Name)
				//}
				//if ok {
				//	// TODO: maybe not take char as input, but generate c.Player with input (character name)
				//	// or possibly this function can return the character to log in to?
				//	go c.enterGame(input)
				//	return nil
				//} else {
				//	c.OutputSteam <- []byte(fmt.Sprintf("You do not have a character named %s", input))
				//	c.state = stateAccountMenu
				//	break
				//}
			}
			break
		case stateAccountMenu:
			c.OutputSteam <- accountMenu
			c.state = statePrompt
			break
		case stateAccountNewCharacter:
			// TODO: character creation
			c.OutputSteam <- "<W>creating character!</W>"
			c.createCharacter()
			break
		case stateAccountListCharacters:
			// TODO: list characters (Account struct method)
			//for _, character in range c.Account.Characters() {
			//	c.OutputSteam <- character
			//}
			// TODO: this might be a deadlock.
			c.OutputSteam <- "<W>Characters: Mrbagginz, Frodo, Samwise</W>"
			c.state = statePrompt
			break
		case stateAccountDeleteCharacter:
			// TODO: delete char prompt
			break
		case stateAccountChangePassword:
			// TODO: change password prompt
			err := c.changePassword()
			if err != nil {
				c.OutputSteam <- err.Error()
			} else {
				c.OutputSteam <- "<Y>Password updated successfully!</Y>"
			}
			c.state = stateAccountMenu
			break
		case stateAccountQuit:
			c.state = stateDisconnected
			return nil
		}
	}
}

func (c *Client) createCharacter() {
	c.OutputSteam <- "<P>You've chosen to create a character!</P>"
	_, err := c.prompt()
	if err != nil {
		fmt.Println("failed to create character")
	}
}

func (c *Client) changePassword() error {
	var password string
	var confirmedPassword string
	var err error
	changeState := stateAccountChangePassword
	for {
		switch c.state {
		case statePrompt:
			switch changeState {
			case stateAccountChangePassword:
				password, err = c.prompt()
				if err != nil {
					return err
				}
				// Reusing statePassword for the confirmation
				c.state = statePassword
				changeState = statePassword
				break
			case statePassword:
				confirmedPassword, err = c.prompt()
				if err != nil {
					return err
				}
				c.state = stateVerify
				break
			}
			break
		case stateAccountChangePassword:
			c.OutputSteam <- "<Y>Enter new password:</Y>"
			c.state = statePrompt
			break
		case statePassword:
			c.OutputSteam <- "<Y>Re-enter to confirm new password:</Y>"
			c.state = statePrompt
			break
		case stateVerify:
			if password != confirmedPassword {
				return fmt.Errorf("password change failed, passwords do not match")
			} else {
				err = c.Account.ChangePassword(password)
				if err != nil {
					return fmt.Errorf("failed to update password: %v", err)
				}
				return nil
			}
		}
	}
}
