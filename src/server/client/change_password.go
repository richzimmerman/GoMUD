package client

import (
	"fmt"
	"utils"
)

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
					return utils.Error(err)
				}
				return nil
			}
		}
	}
}
