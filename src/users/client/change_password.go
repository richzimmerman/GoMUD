package client

import (
	"fmt"
	. "lib/accounts"
	"utils"
)

func (c *Client) changePassword(accountName string) error {
	account, _ := GetAccount(accountName)
	var password string
	var confirmedPassword string
	var err error
	changeState := stateAccountChangePassword
	for {
		switch c.state {
		case statePrompt:
			switch changeState {
			case stateAccountChangePassword:
				password, err = c.Telnet.Prompt()
				if err != nil {
					return utils.Error(err)
				}
				// Reusing statePassword for the confirmation
				c.state = statePassword
				changeState = statePassword
				break
			case statePassword:
				confirmedPassword, err = c.Telnet.Prompt()
				if err != nil {
					return utils.Error(err)
				}
				c.state = stateVerify
				break
			}
			break
		case stateAccountChangePassword:
			c.OutputStream <- "<Y>Enter new password:</Y>"
			c.state = statePrompt
			break
		case statePassword:
			c.OutputStream <- "<Y>Re-enter to confirm new password:</Y>"
			c.state = statePrompt
			break
		case stateVerify:
			if password != confirmedPassword {
				return fmt.Errorf("password change failed, passwords do not match")
			} else {
				err = account.ChangePassword(password)
				if err != nil {
					return utils.Error(err)
				}
				return nil
			}
		}
	}
}
