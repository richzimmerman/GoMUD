package client

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
	stateCreateAccountEmail           = 19
	stateVerify                       = 20

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

func (c *Client) GameLoop() error {
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
					i, err := c.prompt()
					if err != nil {
						return false, err
					}
					if strings.ToLower(i) == "y" {
						return c.createAccountPrompt()
					} else {
						c.state = stateAccount
					}
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

func (c *Client) createAccountPrompt() (bool, error) {
	var accountName string
	var password string
	var email string
	var err error

	c.state = stateCreateAccountName
	var createState = stateCreateAccountName
	for {
		switch c.state {
		case statePrompt:
			switch createState {
			case stateCreateAccount:
				break
			case stateCreateAccountName:
				accountName, err = c.prompt()
				if err != nil {
					return false, err
				}
				c.state = stateCreateAccountNameConfirm
				break
			case stateCreateAccountNameConfirm:
				i, err := c.prompt()
				if err != nil {
					return false, err
				}
				if strings.ToLower(i) == "y" {
					c.state = stateCreateAccountPassword
				} else if strings.ToLower(i) == "n" {
					// Reset account creation (or maybe go back to main menu) (or maybe add Q option)
					createState = stateCreateAccount
					c.state = stateCreateAccount
				} else {
					c.OutputSteam <- "Please enter Y or N to confirm your account name."
				}
				break
			case stateCreateAccountPassword:
				password, err = c.prompt()
				if err != nil {
					return false, err
				}
				// TODO: match password against required password schema
				c.state = stateCreateAccountPasswordConfirm
				break
			case stateCreateAccountPasswordConfirm:
				i, err := c.prompt()
				if err != nil {
					return false, err
				}
				if i != password {
					c.OutputSteam <- "<Y>Password does not match, please retry.</Y>"
				}
				c.state = stateCreateAccountEmail
				break
			case stateCreateAccountEmail:
				email, err = c.prompt()
				if err != nil {
					return false, err
				}
				c.state = stateCreateAccount
				break
			}
			break
		case stateCreateAccountName:
			c.OutputSteam <- "<Y>Enter desired account name (Note: Account names are case sensitive):</Y>"
			c.state = statePrompt
			break
		case stateCreateAccountNameConfirm:
			// TODO: db verify account name
			// if account name exists, error back out to main menu
			c.OutputSteam <- fmt.Sprintf("%s: <Y>is this correct? (Y/n)</Y>", accountName)
			c.state = statePrompt
			createState = stateCreateAccountNameConfirm
			break
		case stateCreateAccountPassword:
			c.OutputSteam <- "<Y>Enter your password:</Y>"
			c.state = statePrompt
			createState = stateCreateAccountPassword
			break
		case stateCreateAccountPasswordConfirm:
			c.OutputSteam <- "<Y>Re-enter password to confirm.</Y>"
			c.state = statePrompt
			createState = stateCreateAccountPasswordConfirm
			break
		case stateCreateAccountEmail:
			c.OutputSteam <- "<Y>Enter your e-mail address.</Y>"
			c.state = statePrompt
			createState = stateCreateAccountEmail
			break
		case stateCreateAccount:
			lastip := strings.Split(c.Connection.RemoteAddr().String(), ":")[0]
			c.Account, err = account.NewAccount(accountName, password, lastip, email)
			if err != nil {
				return false, err
			}

			return true, nil
		}
	}
}

func (c *Client) createAccount(accountName string, password string, email string) error {
	lastip := strings.SplitN(c.Connection.RemoteAddr().String(), ":", 1)
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
			c.OutputSteam <- "creating character!"
			c.createCharacter()
			break
		case stateAccountListCharacters:
			// TODO: list characters (Account struct method)
			if len(c.Account.Characters) == 0 {
				c.OutputSteam <- "<Y>You have no characters on this account.</Y>"
			} else {
				for _, character := range c.Account.Characters {
					c.OutputSteam <- fmt.Sprintf("<Y>%s</Y>: Level %d %s",
						character.Name(), character.Level(), character.Race())
				}
			}
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
