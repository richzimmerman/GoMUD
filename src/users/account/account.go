package account

import (
	"db"
	. "db"
	"fmt"
	"interfaces"
	"logger"
	"utils"
)

var log = logger.NewLogger()

type Account struct {
	Characters []string
	LoggedIn   bool
	*DBAccount
}

func LoadAccount(accountName string) (*Account, error) {
	/*
		Loads the account info at login time, also loads up all of the characters into the account Character map
	*/
	dba, err := db.DatabaseConnection.LoadAccount(accountName)
	if err != nil {
		return nil, err
	}
	chars, err := db.DatabaseConnection.LoadAccountCharacters(accountName)
	if err != nil {
		return nil, fmt.Errorf("failed to load players when loading account: %v", err)
	}
	characters := make([]string, 0)
	for name, _ := range chars {
		characters = append(characters, name)
	}
	a := &Account{
		Characters: characters,
		LoggedIn:   true,
		DBAccount:  dba,
	}
	return a, nil
}

func (a *Account) AccountName() string {
	return a.Name
}

func (a *Account) LastIPLogged() string {
	return a.LastIP
}

func (a *Account) EmailAddress() string {
	return a.Email
}

func (a *Account) GetCharacters() []string {
	return a.Characters
}

func (a *Account) GetPassword() string {
	return a.Password
}

func (a *Account) LoggedInStatus() bool {
	return a.LoggedIn
}

func (a *Account) SetLoggedInStatus(b bool) {
	a.LoggedIn = b
}

func (a *Account) CreateCharacter(p interfaces.PlayerInterface) error {
	// TODO: delete this i think
	// TODO: use db to create characters for this account, add player to list of players
	fmt.Printf("Creating character: %s on account %s\n", p.GetName(), a.Name)
	return nil
}

func (a *Account) VerifyCharacter(characterName string) bool {
	/*
		Helper function to verify character exists. Characters are loaded upon login so this just checks the map of
		Characters.
	*/
	return utils.ContainsString(a.Characters, characterName)
}

func (a *Account) DeleteCharacter(characterName string) (bool, error) {
	// TODO: use db to delete character, make sure to delete character from a.Characters
	fmt.Printf("Deleting character: %s \n", characterName)
	i := utils.IndexOf(characterName, a.Characters)
	if i < 0 {
		return false, fmt.Errorf("unable to delete character: %s", characterName)
	}
	copy(a.Characters[i:], a.Characters[i+1:])        // Shift a[i+1:] left one index.
	a.Characters[len(a.Characters)-1] = ""            // Erase last element (write zero value).
	a.Characters = a.Characters[:len(a.Characters)-1] // Truncate slice.
	return true, nil
}

func (a *Account) ChangePassword(s string) error {
	// TODO: match input against required password schema.
	err := db.DatabaseConnection.ChangePassword(a.Name, s)
	if err != nil {
		return fmt.Errorf("unable to change password for account %s: %v", a.Name, err)
	}
	a.Password = s
	return nil
}
