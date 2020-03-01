package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strings"
)

var DatabaseConnection *dbConnection

// TODO: Signing key for password encryption/decryption needs to be Configurable
const key = "signing_key"

type dbConnection struct {
	db *sql.DB
}

type DBAccount struct {
	Name     string
	Password string
	LastIP   string
	Email    string
}

type DBPlayer struct {
	Name        string
	Account     string
	DisplayName string
	Level       int8
	Health      int16
	Fatigue     int16
	Power       int16
	Title       string
	RealmTitle  string
	Race        string
	Stats       string // JSON?
	Stance      int8
	Skills      string // JSON Array: List of skills to load later ["skill1", "skill2", "skill3"]
	Spells      string // Same as above
	Buffs       string // JSON array of objects [{"name": "buff1", "duration":12345}, {}, {}]
	Debuffs     string // Same as above
}

func InitDatabaseConnection() error {
	db, err := sql.Open("mysql", "gomud:test@/GoMUD")
	if err != nil {
		return err
	}
	DatabaseConnection = &dbConnection{
		db: db,
	}
	return nil
}

func (d *dbConnection) CreateAccount(a *DBAccount) error {
	// TODO: Test this
	s := fmt.Sprintf("INSERT INTO Accounts VALUES (?, AES_ENCRYPT(?, '%s'), ?, ?)", key)
	statement, err := d.db.Prepare(s)
	if err != nil {
		return fmt.Errorf("unable to create account: %v", err)
	}
	defer statement.Close()

	res, err := statement.Exec(a.Name, a.Password, a.LastIP, a.Email)
	if err != nil {
		return fmt.Errorf("error executing insert account statement: %v", err)
	}
	_, err = res.RowsAffected()
	if err != nil {
		return fmt.Errorf("unable to get rows inserted: %v", err)
	}
	return nil
}

func (d *dbConnection) AccountExists(accountName string) (bool, error) {
	searchStatement, err := d.db.Prepare("SELECT Name FROM Accounts WHERE Name = ?")
	if err != nil {
		return false, err
	}
	defer searchStatement.Close()

	var a string
	err = searchStatement.QueryRow(accountName).Scan(&a)
	if err != nil {
		// This should indicate that now row was returned and the account does not exist.
		fmt.Println(fmt.Errorf("unable to get account id: %v", err))
		if strings.Contains(err.Error(), "no rows in result") {
			return false, nil
		} else {
			return false, err
		}
	}

	// For sanity's sake.
	if a != "" {
		return true, nil
	} else {
		return false, nil
	}
}

func (d *dbConnection) VerifyPassword(accountName string, password string) (bool, error) {
	statement := fmt.Sprintf("SELECT AES_DECRYPT(Password, '%s') FROM Accounts WHERE Name = ?", key)
	rows, err := d.db.Query(statement, accountName)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	var pw string
	for rows.Next() {
		err := rows.Scan(&pw)
		if err != nil {
			return false, err
		}
	}
	return pw == password, nil
}

func (d *dbConnection) LoadAccount(accountName string) (*DBAccount, error) {
	searchStatement, err := d.db.Prepare("SELECT * FROM Accounts WHERE Name = ?")
	if err != nil {
		return nil, err
	}
	defer searchStatement.Close()

	a := &DBAccount{}
	err = searchStatement.QueryRow(accountName).Scan(&a.Name, &a.Password, &a.LastIP, &a.Email)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (d *dbConnection) LoadAccountCharacters(accountName string) (map[string]*DBPlayer, error) {
	// TODO: verify this works right
	r := make(map[string]*DBPlayer)
	rows, err := d.db.Query("SELECT * FROM Characters WHERE Account = ?", accountName)
	if err != nil {
		return nil, fmt.Errorf("unable to query characters for account %s: %v", accountName, err)
	}
	defer rows.Close()

	for rows.Next() {
		p := &DBPlayer{}
		err := rows.Scan(&p.Name, &p.Account, &p.DisplayName, &p.Level, &p.Health, &p.Fatigue, &p.Power, &p.Title,
			&p.RealmTitle, &p.Race, &p.Stats, &p.Stance, &p.Skills, &p.Spells, &p.Buffs, &p.Debuffs)
		if err != nil {
			fmt.Printf("failed to scan row: %v\n", err)
			return nil, err
		}
		r[p.Name] = p
	}
	return r, nil
}

func (d *dbConnection) ChangePassword(accountName string, password string) error {
	s := fmt.Sprintf("UPDATE Accounts SET Password = AES_ENCRYPT(?, '%s') WHERE Name = ?", key)
	statement, err := d.db.Prepare(s)
	if err != nil {
		return fmt.Errorf("unable to create account: %v", err)
	}
	defer statement.Close()

	res, err := statement.Exec(password, accountName)
	if err != nil {
		return fmt.Errorf("error executing update password statement: %v", err)
	}
	_, err = res.RowsAffected()
	if err != nil {
		return fmt.Errorf("unable to get rows updated: %v", err)
	}
	return nil
}
