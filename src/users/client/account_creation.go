package client

import (
	"fmt"
	"lib/accounts"
	"strings"
	"utils"
)

func (c *Client) createAccountPrompt() (string, error) {
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
				accountName, err = c.Telnet.Prompt()
				if err != nil {
					return "", utils.Error(err)
				}
				c.state = stateCreateAccountNameConfirm
				break
			case stateCreateAccountNameConfirm:
				i, err := c.Telnet.Prompt()
				if err != nil {
					return "", utils.Error(err)
				}
				if strings.ToLower(i) == "y" {
					c.SetAssociatedAccount(accountName)
					c.state = stateCreateAccountPassword
				} else if strings.ToLower(i) == "n" {
					// Reset account creation (or maybe go back to main menu) (or maybe add Q option)
					createState = stateCreateAccount
					c.state = stateCreateAccount
				} else {
					c.OutputStream <- "Please enter Y or N to confirm your account name."
				}
				break
			case stateCreateAccountPassword:
				password, err = c.Telnet.Prompt()
				if err != nil {
					return "", utils.Error(err)
				}
				// TODO: match password against required password schema
				c.state = stateCreateAccountPasswordConfirm
				break
			case stateCreateAccountPasswordConfirm:
				i, err := c.Telnet.Prompt()
				if err != nil {
					return "", utils.Error(err)
				}
				if i != password {
					c.OutputStream <- "<Y>Password does not match, please retry.</Y>"
				}
				c.state = stateCreateAccountEmail
				break
			case stateCreateAccountEmail:
				email, err = c.Telnet.Prompt()
				if err != nil {
					return "", utils.Error(err)
				}
				c.state = stateCreateAccount
				break
			}
			break
		case stateCreateAccountName:
			c.OutputStream <- "<Y>Enter desired account name (Note: Account names are case sensitive):</Y>"
			c.state = statePrompt
			break
		case stateCreateAccountNameConfirm:
			// TODO: db verify account name
			// if account name exists, error back out to main menu
			c.OutputStream <- fmt.Sprintf("%s: <Y>is this correct? (Y/n)</Y>", accountName)
			c.state = statePrompt
			createState = stateCreateAccountNameConfirm
			break
		case stateCreateAccountPassword:
			c.OutputStream <- "<Y>Enter your password:</Y>"
			c.state = statePrompt
			createState = stateCreateAccountPassword
			break
		case stateCreateAccountPasswordConfirm:
			c.OutputStream <- "<Y>Re-enter password to confirm.</Y>"
			c.state = statePrompt
			createState = stateCreateAccountPasswordConfirm
			break
		case stateCreateAccountEmail:
			c.OutputStream <- "<Y>Enter your e-mail address.</Y>"
			c.state = statePrompt
			createState = stateCreateAccountEmail
			break
		case stateCreateAccount:
			lastip := strings.Split(c.Connection.RemoteAddr().String(), ":")[0]
			_, err = accounts.NewAccount(accountName, password, lastip, email)
			if err != nil {
				return "", utils.Error(err)
			}
			return accountName, nil
		}
	}
}
