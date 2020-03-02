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
			input, err = c.prompt()
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
			c.OutputSteam <- "Disconnected!"
			c.state = stateDisconnected
			return nil
		}
	}
}
