package client

import (
	"fmt"
	. "lib/accounts"
	. "lib/players"
	"strings"
	"utils"
)

func (c *Client) accountMenu(accountName string) error {
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
				player, err := CheckPlayer(input, nil)
				if err != nil || player.AccountName() != accountName {
					c.OutputStream <- "<Y>The character does not exist!</Y>"
				} else {
					c.SetAssociatedPlayer(input)
					c.state = stateInGame
					return nil
				}
			}
			break
		case stateAccountMenu:
			c.OutputStream <- accountMenu
			c.state = statePrompt
			break
		case stateAccountNewCharacter:
			err := c.createCharacter(accountName)
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
			acct, err := GetAccount(accountName)
			if err != nil {
				log.Info("uhhh wtf")
			}
			if len(acct.GetCharacters()) == 0 {
				c.OutputStream <- "<Y>You have no characters on this account.</Y>"
			} else {
				for _, character := range acct.GetCharacters() {
					p, err := CheckPlayer(character, nil)
					if err != nil {
						log.Err("account (%s) has a character (%s) that does not exist", acct.AccountName(), character)
					}
					c.OutputStream <- fmt.Sprintf("<Y>%s</Y>: Level %d %s",
						p.GetName(), p.GetLevel(), p.RaceName())
				}
			}
			c.state = statePrompt
			break
		case stateAccountDeleteCharacter:
			// TODO: delete char prompt
			break
		case stateAccountChangePassword:
			// TODO: change password prompt
			err := c.changePassword(accountName)
			if err != nil {
				c.OutputStream <- err.Error()
			} else {
				c.OutputStream <- "<Y>Password updated successfully!</Y>"
			}
			c.state = stateAccountMenu
			break
		case stateAccountQuit:
			c.OutputStream <- "Disconnected!"
			c.state = stateDisconnected
			return nil
		}
	}
}
