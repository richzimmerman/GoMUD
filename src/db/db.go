package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strings"
	"utils"
)

var DatabaseConnection *DbConnection

// TODO: Signing key for password encryption/decryption needs to be Configurable
const key = "signing_key"

type DbConnection struct {
	Connection *sql.DB
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

type DBRace struct {
	Name        string
	Realm       int8
	Type        int8
	SkillList   string // JSON array
	Description string
}

func InitDatabaseConnection() error {
	db, err := sql.Open("mysql", "gomud:test@/GoMUD")
	if err != nil {
		return utils.Error(err)
	}
	DatabaseConnection = &DbConnection{
		Connection: db,
	}
	return nil
}

func (d *DbConnection) CreateAccount(a *DBAccount) error {
	s := fmt.Sprintf("INSERT INTO Accounts VALUES (?, AES_ENCRYPT(?, '%s'), ?, ?)", key)
	statement, err := d.Connection.Prepare(s)
	if err != nil {
		return utils.Error(err)
	}
	defer statement.Close()

	res, err := statement.Exec(a.Name, a.Password, a.LastIP, a.Email)
	if err != nil {
		return utils.Error(err)
	}
	_, err = res.RowsAffected()
	if err != nil {
		return utils.Error(err)
	}
	return nil
}

func (d *DbConnection) AccountExists(accountName string) (bool, error) {
	searchStatement, err := d.Connection.Prepare("SELECT Name FROM Accounts WHERE Name = ?")
	if err != nil {
		return false, utils.Error(err)
	}
	defer searchStatement.Close()

	var a string
	err = searchStatement.QueryRow(accountName).Scan(&a)
	if err != nil {
		// This should indicate that now row was returned and the account does not exist.
		if strings.Contains(err.Error(), "no rows in result") {
			return false, nil
		} else {
			return false, utils.Error(err)
		}
	}

	// For sanity's sake.
	if a != "" {
		return true, nil
	} else {
		return false, nil
	}
}

func (d *DbConnection) VerifyPassword(accountName string, password string) (bool, error) {
	statement := fmt.Sprintf("SELECT AES_DECRYPT(Password, '%s') FROM Accounts WHERE Name = ?", key)
	rows, err := d.Connection.Query(statement, accountName)
	if err != nil {
		return false, utils.Error(err)
	}
	defer rows.Close()

	var pw string
	for rows.Next() {
		err := rows.Scan(&pw)
		if err != nil {
			return false, utils.Error(err)
		}
	}
	return pw == password, nil
}

func (d *DbConnection) LoadAccount(accountName string) (*DBAccount, error) {
	searchStatement, err := d.Connection.Prepare("SELECT * FROM Accounts WHERE Name = ?")
	if err != nil {
		return nil, utils.Error(err)
	}
	defer searchStatement.Close()

	a := &DBAccount{}
	err = searchStatement.QueryRow(accountName).Scan(&a.Name, &a.Password, &a.LastIP, &a.Email)
	if err != nil {
		return nil, utils.Error(err)
	}
	return a, nil
}

func (d *DbConnection) LoadAccountCharacters(accountName string) (map[string]*DBPlayer, error) {
	r := make(map[string]*DBPlayer)
	rows, err := d.Connection.Query("SELECT * FROM Characters WHERE Account = ?", accountName)
	if err != nil {
		return nil, utils.Error(err)
	}
	defer rows.Close()

	for rows.Next() {
		p := &DBPlayer{}
		err := rows.Scan(&p.Name, &p.Account, &p.DisplayName, &p.Level, &p.Health, &p.Fatigue, &p.Power, &p.Title,
			&p.RealmTitle, &p.Race, &p.Stats, &p.Stance, &p.Skills, &p.Spells, &p.Buffs, &p.Debuffs)
		if err != nil {
			return nil, utils.Error(err)
		}
		r[p.Name] = p
	}
	return r, nil
}

func (d *DbConnection) ChangePassword(accountName string, password string) error {
	s := fmt.Sprintf("UPDATE Accounts SET Password = AES_ENCRYPT(?, '%s') WHERE Name = ?", key)
	statement, err := d.Connection.Prepare(s)
	if err != nil {
		return utils.Error(err)
	}
	defer statement.Close()

	res, err := statement.Exec(password, accountName)
	if err != nil {
		return utils.Error(err)
	}
	_, err = res.RowsAffected()
	if err != nil {
		return utils.Error(err)
	}
	return nil
}

func (d *DbConnection) CreateRace(r *DBRace) error {
	// TODO: Test this (in game)
	statement, err := d.Connection.Prepare("INSERT INTO Races VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return utils.Error(err)
	}
	defer statement.Close()

	res, err := statement.Exec(r.Name, r.Realm, r.Type, r.SkillList, r.Description)
	if err != nil {
		return utils.Error(err)
	}
	_, err = res.RowsAffected()
	if err != nil {
		return utils.Error(err)
	}
	return nil
}

func (d *DbConnection) LoadRaces() ([]*DBRace, error) {
	var races []*DBRace
	rows, err := d.Connection.Query("SELECT * FROM Races")
	if err != nil {
		return nil, utils.Error(err)
	}
	defer rows.Close()

	for rows.Next() {
		r := &DBRace{}
		err := rows.Scan(&r.Name, &r.Realm, &r.Type, &r.SkillList, &r.Description)
		if err != nil {
			return nil, utils.Error(err)
		}
		races = append(races, r)
	}
	return races, nil
}
