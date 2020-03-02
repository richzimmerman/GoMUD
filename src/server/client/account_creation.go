package client

import (
	"account"
	"fmt"
	"strings"
	"utils"
)

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
					return false, utils.Error(err)
				}
				c.state = stateCreateAccountNameConfirm
				break
			case stateCreateAccountNameConfirm:
				i, err := c.prompt()
				if err != nil {
					return false, utils.Error(err)
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
					return false, utils.Error(err)
				}
				// TODO: match password against required password schema
				c.state = stateCreateAccountPasswordConfirm
				break
			case stateCreateAccountPasswordConfirm:
				i, err := c.prompt()
				if err != nil {
					return false, utils.Error(err)
				}
				if i != password {
					c.OutputSteam <- "<Y>Password does not match, please retry.</Y>"
				}
				c.state = stateCreateAccountEmail
				break
			case stateCreateAccountEmail:
				email, err = c.prompt()
				if err != nil {
					return false, utils.Error(err)
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
				return false, utils.Error(err)
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
		return utils.Error(err)
	}
	c.Account = a
	return nil
}