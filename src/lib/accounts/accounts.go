package accounts

import (
	"db"
	"fmt"
	. "interfaces"
	"logger"
	"users/account"
)

var log = logger.NewLogger()

var accounts = make(map[string]AccountInterface)

func LoadAccounts() error {
	allAccounts, err := db.DatabaseConnection.LoadAllAccounts()
	if err != nil {
		return err
	}

	for acctName, acct := range allAccounts {
		chars, err := db.DatabaseConnection.LoadAccountCharacters(acctName)
		if err != nil {
			return fmt.Errorf("failed to load players when loading account: %v", err)
		}
		characters := make([]string, 0)
		for name, _ := range chars {
			characters = append(characters, name)
		}
		a := &account.Account{
			Characters: characters,
			DBAccount:  acct,
		}
		accounts[acctName] = a
	}
	return nil
}

func AddAccount(a AccountInterface) error {
	if _, found := accounts[a.AccountName()]; found {
		return fmt.Errorf("account (%s) already exists", a.AccountName())
	}
	accounts[a.AccountName()] = a
	return nil
}

func GetAccount(name string) (AccountInterface, error) {
	if a, found := accounts[name]; found {
		log.Info("returning account: %s", a.AccountName())
		return a, nil
	}
	return nil, fmt.Errorf("account (%s) not found", name)
}

func RemoveAccount(name string) error {
	if _, found := accounts[name]; !found {
		return fmt.Errorf("account (%s) not found", name)
	}
	delete(accounts, name)
	return nil
}

func NewAccount(accountName string, password string, lastip string, email string) (AccountInterface, error) {
	dba := &db.DBAccount{
		Name:     accountName,
		Password: password,
		LastIP:   lastip,
		Email:    email,
	}
	a := &account.Account{
		DBAccount:  dba,
		Characters: make([]string, 0),
	}
	err := db.DatabaseConnection.CreateAccount(dba)
	if err != nil {
		return nil, err
	}
	AddAccount(a)
	return a, nil
}
