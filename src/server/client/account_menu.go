package client

import (
	"fmt"
	"strings"
	"utils"
)

func (c *Client) accountMenu() error {
	var input string
	var err error
	c.state = stateAccountMenu
	for {
		switch c.state {
		case statePrompt:
			input, err = c.Telnet.Prompt()
			if err != nil {
				return utils.Error(err)
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
				input = strings.Title(input)
				player, ok := c.Account.Characters[input]
				if !ok {
					c.OutputSteam <- "<Y>The character does not exist!</Y>"
				} else {
					c.Player = player
					c.state = stateInGame
					return nil
				}
			}
			break
		case stateAccountMenu:
			c.OutputSteam <- accountMenu
			c.state = statePrompt
			break
		case stateAccountNewCharacter:
			err := c.createCharacter()
			if err != nil {
				return utils.Error(err)
			}
			if c.state == stateAccountMenu {
				break
			} else {
				return nil
			}
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
			c.OutputSteam <- "Disconnected!"
			c.state = stateDisconnected
			return nil
		}
	}
}
