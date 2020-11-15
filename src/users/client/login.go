package client

import (
	"db"
	. "lib/accounts"
	"strings"
	"utils"
)

func (c *Client) logIn() (string, error) {
	var accountName string
	var err error
	loginState := stateAccount
	c.state = stateAccount
	for {
		switch c.state {
		case statePrompt:
			switch loginState {
			case stateAccount:
				accountName, err = c.Telnet.Prompt()
				if err != nil {
					return "", utils.Error(err)
				}
				ok, err := db.DatabaseConnection.AccountExists(accountName)
				if err != nil {
					return "", utils.Error(err)
				}
				if ok {
					loginState = statePassword
					c.state = statePassword
				} else {
					c.OutputStream <- invalidAccount
					i, err := c.Telnet.Prompt()
					if err != nil {
						return "", utils.Error(err)
					}
					if strings.ToLower(i) == "y" {
						return c.createAccountPrompt()
					} else {
						c.state = stateAccount
					}
				}
				break
			case statePassword:
				pw, err := c.Telnet.Prompt()
				if err != nil {
					return "", err
				}
				ok, err := db.DatabaseConnection.VerifyPassword(accountName, pw)
				if err != nil {
					return "", err
				}
				if ok {
					a, err := GetAccount(accountName)
					if err != nil {
						return "", utils.Error(err)
					}
					if a.LoggedInStatus() {
						c.OutputStream <- accountAlreadyLoggedIn
						loginState = stateAccount
						c.state = stateAccount
						break
					}
					c.SetAssociatedAccount(accountName)
					a.SetLoggedInStatus(true)
					return a.AccountName(), nil
				} else {
					c.OutputStream <- invalidPassword
					loginState = stateAccount
					c.state = stateAccount
				}
				break
			}
			break
		case stateAccount:
			c.OutputStream <- accountPrompt
			c.state = statePrompt
			break
		case statePassword:
			c.OutputStream <- passwordPrompt
			c.state = statePrompt
			break
		}
	}
}
