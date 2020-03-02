package client

import (
	"account"
	"db"
	"strings"
)

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
