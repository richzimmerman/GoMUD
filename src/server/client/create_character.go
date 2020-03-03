package client

import (
	"db"
	"fmt"
	"mobs"
	"races"
	"strings"
	"utils"
	"world/realms"
)

const (
	stateName = 1
	stateNameConfirm = 2
	stateRealm = 3
	stateRealmConfirm = 4
	stateRace = 5
	stateRaceConfirm = 6
	stateCharacterConfirm = 7
)

func generateRaceOutput(realm int8) (string, error) {
	output := "Please select one of the following races:\n"
	for name, race := range races.Races {
		if race.Realm == realm {
			output = output + name + "\n"
		}
	}
	if output == "Please select one of the following races:\n" {
		return "", utils.Error(fmt.Errorf("no races available for realm %d", realm))
	}
	return output, nil
}

func (c *Client) createCharacter() error {
	var name string
	var realm string
	var chosenRealm int8
	var race string
	var err error

	createState := stateName
	c.state = stateName
	for {
		switch c.state {
		case statePrompt:
			switch createState {
			case stateName:
				name, err = c.Telnet.Prompt()
				if err != nil {
					return utils.Error(err)
				}
				if name == "" {
					c.OutputSteam <- "<Y>Please enter a name.</Y>"
					break
				}
				if quit := utils.CheckIfQuit(name); quit {
					c.state = stateAccountMenu
					return nil
				}
				name = strings.Title(name)
				ok, err := db.DatabaseConnection.CharacterNameAvailable(name)
				if err != nil {
					return utils.Error(err)
				}
				if !ok {
					c.OutputSteam <- fmt.Sprintf("<Y>%s is not available, please choose another name.</Y>", name)
				} else {
					c.state = stateNameConfirm
				}
				break
			case stateNameConfirm:
				input, err := c.Telnet.Prompt()
				if err != nil {
					return utils.Error(err)
				}
				if quit := utils.CheckIfQuit(input); quit {
					c.state = stateAccountMenu
					return nil
				}
				switch strings.ToLower(input) {
				case "y":
					c.state = stateRealm
					break
				case "n":
					c.state = stateName
					break
				default:
					c.OutputSteam <- "<Y>Please enter Y or N to confirm</Y>"
				}
			case stateRealm:
				realm, err = c.Telnet.Prompt()
				if err != nil {
					return utils.Error(err)
				}
				if realm == "" {
					c.OutputSteam <- "<Y>Please enter a realm.</Y>"
					break
				}
				if quit := utils.CheckIfQuit(realm); quit {
					c.state = stateAccountMenu
					return nil
				}
				realm = strings.Title(realm)
				chosenRealm = int8(utils.IndexOf(realm, realms.Realms))
				if chosenRealm == -1 {
					c.OutputSteam <- "<Y>Please choose a realm</Y>"
				} else {
					c.state = stateRealmConfirm
				}
				break
			case stateRealmConfirm:
				input, err := c.Telnet.Prompt()
				if err != nil {
					return utils.Error(err)
				}
				if quit := utils.CheckIfQuit(input); quit {
					c.state = stateAccountMenu
					return nil
				}
				switch strings.ToLower(input) {
				case "y":
					c.state = stateRace
					break
				case "n":
					c.state = stateRealm
					break
				default:
					c.OutputSteam <- "<Y>Please enter Y or N to confirm</Y>"
				}
			case stateRace:
				race, err = c.Telnet.Prompt()
				if err != nil {
					return utils.Error(err)
				}
				if race == "" {
					c.OutputSteam <- "<Y>Please enter a race.</Y>"
					break
				}
				if quit := utils.CheckIfQuit(race); quit {
					c.state = stateAccountMenu
					return nil
				}
				race = strings.Title(race)
				_, ok := races.Races[race]
				if !ok {
					c.OutputSteam <- "<Y>Please select a valid race.</Y>"
				} else {
					c.state = stateRaceConfirm
				}
				break
			case stateRaceConfirm:
				input, err := c.Telnet.Prompt()
				if err != nil {
					return utils.Error(err)
				}
				if quit := utils.CheckIfQuit(input); quit {
					c.state = stateAccountMenu
					return nil
				}
				switch strings.ToLower(input) {
				case "y":
					c.state = stateCharacterConfirm
					break
				case "n":
					c.state = stateRace
					break
				default:
					c.OutputSteam <- "<Y>Please enter Y or N to confirm</Y>"
				}
			case stateCharacterConfirm:
				input, err := c.Telnet.Prompt()
				if err != nil {
					return utils.Error(err)
				}
				if quit := utils.CheckIfQuit(input); quit {
					c.state = stateAccountMenu
					return nil
				}
				switch strings.ToLower(input) {
				case "y":
					c.state = stateInGame
					p, err := mobs.NewPlayer(c.Account.Name, name, race, chosenRealm)
					if err != nil {
						return utils.Error(err)
					}
					c.Player = p
					return nil
				case "n":
					c.state = stateRace
					break
				default:
					c.OutputSteam <- "<Y>Please enter Y or N to confirm</Y>"
				}
			}
			break
		case stateName:
			c.OutputSteam <- "<Y>Enter your new character's name.</Y>"
			c.state = statePrompt
			createState = stateName
			break
		case stateNameConfirm:
			c.OutputSteam <- fmt.Sprintf("<Y>%s</Y>, is this correct? <Y>(Y/n)</Y>", name)
			c.state = statePrompt
			createState = stateNameConfirm
			break
		case stateRealm:
			c.OutputSteam <- "Choose your realm: <Y>Good</Y>, <Y>Chaos</Y>, or <Y>Evil</Y>"
			c.state = statePrompt
			createState = stateRealm
			break
		case stateRealmConfirm:
			c.OutputSteam <- fmt.Sprintf("<Y>%s</Y>, is this correct? <Y>(Y/n)</Y>", realm)
			c.state = statePrompt
			createState = stateRealmConfirm
			break
		case stateRace:
			output, err := generateRaceOutput(chosenRealm)
			if err != nil {
				return utils.Error(err)
			}
			c.OutputSteam <- output
			c.state = statePrompt
			createState = stateRace
			break
		case stateRaceConfirm:
			c.OutputSteam <- fmt.Sprintf("<Y>%s</Y>, is this correct? <Y>(Y/n)</Y>", race)
			c.state = statePrompt
			createState = stateRaceConfirm
			break
		case stateCharacterConfirm:
			s := fmt.Sprintf("<Y>You've chosen %s, a %s in the %s realm. Is this correct?</Y>", name, race, realm)
			c.OutputSteam <- s
			c.state = statePrompt
			createState = stateCharacterConfirm
			break
		}
	}
}