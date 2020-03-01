package account

import (
	"db"
	"fmt"
	"mobs"
)

type Account struct {
	*db.DBAccount
	Characters map[string]*mobs.Player
}

func NewAccount(accountName string, password string, lastip string, email string) (*Account, error) {
	dba := &db.DBAccount{
		Name:     accountName,
		Password: password,
		LastIP:   lastip,
		Email:    email,
	}
	a := &Account{
		dba,
		make(map[string]*mobs.Player),
	}
	err := db.DatabaseConnection.CreateAccount(dba)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func LoadAccount(accountName string) (*Account, error) {
	/*
	Loads the account info at login time, also loads up all of the characters into the account Character map
	 */
	characters := make(map[string]*mobs.Player)
	dba, err := db.DatabaseConnection.LoadAccount(accountName)
	if err != nil {
		return nil, err
	}
	chars, err := db.DatabaseConnection.LoadAccountCharacters(accountName)
	if err != nil {
		return nil, fmt.Errorf("failed to load players when loading account: %v", err)
	}
	for name, character := range chars {
		p, err := mobs.LoadPlayer(character)
		if err != nil {
			return nil, err
		}
		characters[name] = p
	}
	return &Account{
		dba,
		characters,
	}, nil
}

func (a *Account) CreateCharacter(p *mobs.Player) error {
	// TODO: use db to create characters for this account, add player to list of players
	fmt.Printf("Creating character: %s on account %s\n", p.Name(), a.Name)
	return nil
}

func (a *Account) VerifyCharacter(characterName string) bool {
	/*
	Helper function to verify character exists. Characters are loaded upon login so this just checks the map of
	Characters.
	 */
	_, ok := a.Characters[characterName]
	if !ok {
		return false
	}
	return true
}

func (a *Account) DeleteCharacter(characterName string) (bool, error) {
	// TODO: use db to delete character, make sure to delete character from a.Characters
	fmt.Printf("Deleting character: %s \n", characterName)
	delete(a.Characters, characterName)
	_, ok := a.Characters[characterName]
	if ok {
		return false, fmt.Errorf("unable to delete character: %s", characterName)
	}
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