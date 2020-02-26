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
	// TODO: make sure to load Characters as well.
	characters := make(map[string]*mobs.Player)
	dba, err := db.DatabaseConnection.LoadAccount(accountName)
	if err != nil {
		return nil, err
	}
	chars, err := db.DatabaseConnection.LoadPlayers(accountName)
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

func (a *Account) loadCharacters() error {
	// TODO: Load characters (Player struct) for a given account from DB
	//[]players := db.DatabaseConnection.LoadPlayers(a.Name)
	//a.Characters[p.Name()] = p
	fmt.Printf("loading characters for account: %s\n", a.Name)
	return nil
}

func (a *Account) CreateCharacter(p *mobs.Player) error {
	// TODO: use db to create characters for this account, add player to list of players
	fmt.Printf("Creating character: %s on account %s\n", p.Name(), a.Name)
	return nil
}

func (a *Account) VerifyCharacter(characterName string) (bool, error) {
	// TODO: match against characters map
	return true, nil
}

func (a *Account) DeleteCharacter(characterName string) (bool, error) {
	// TODO: use db to delete character, make sure to delete character from a.Characters
	fmt.Printf("Deleting character: %s \n", characterName)
	return true, nil
}

func (a *Account) ChangePassword(s string) error {
	err := db.DatabaseConnection.ChangePassword(a.Name, s)
	if err != nil {
		return fmt.Errorf("unable to change password for account %s: %v", a.Name, err)
	}
	a.Password = s
	return nil
}