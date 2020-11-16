package client

import (
	"db"
	"fmt"
	"lib/accounts"
	. "lib/players"
	"net"
	"races"
	"sync"
	"testing"
	. "tests"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func callAccountMenu(client *Client, wg *sync.WaitGroup) error {
	defer wg.Done()
	err := client.accountMenu("TestAccount")
	return err
}

func TestAccountMenu(t *testing.T) {
	var err error
	var conn net.Conn

	mock, err := InitMockDB()
	if err != nil {
		t.Fatalf("error mocking database: %v", err)
	}
	defer db.DatabaseConnection.Connection.Close()

	columnNames := []string{
		"Name", "Account", "DisplayName", "Level", "Health", "Fatigue", "Power", "Title", "RealmTitle", "Race", "Stats",
		"Stance", "Skills", "Spells", "Buffs", "Debuffs", "Location",
	}
	characterQuery := sqlmock.NewRows(columnNames).AddRow("TestPlayer", "TestAccount", "TestPlayer", "1", "1", "1",
		"1", "SlayerBro", "Slayer of Bros", "TestClass", "{\"strength\": 20, \"agility\": 20}",
		"0", "[]", "[]", "{}", "{}", "0")

	mock.ExpectPrepare("SELECT \\* FROM Characters").ExpectQuery().WithArgs("TestPlayer").WillReturnRows(characterQuery)

	acct := NewMockAccount()
	accounts.AddAccount(acct)
	defer accounts.RemoveAccount(acct.AccountName())

	races.Races = make(map[string]*races.Race)

	races.Races["TestClass"] = &races.Race{
		Name:           "TestClass",
		Realm:          0,
		Type:           0,
		SkillList:      nil,
		Description:    "",
		DefaultHealth:  0,
		DefaultFatigue: 0,
		DefaultPower:   0,
		StartingRoom:   "",
		DefaultTitle:   "",
		DefaultStats:   make(map[string]int8),
	}

	p := &db.DBPlayer{
		Name:        "TestPlayer",
		Account:     "TestAccount",
		DisplayName: "TestPlayer",
		Level:       1,
		Health:      100,
		Fatigue:     100,
		Power:       100,
		Title:       "TestTitle",
		RealmTitle:  "TestRealmTitle",
		Race:        "TestClass",
		Stats:       "{}",
		Stance:      0,
		Skills:      "[]",
		Spells:      "[]",
		Buffs:       "{}",
		Debuffs:     "{}",
	}

	l, err := net.Listen("tcp4", ":54321")
	if err != nil {
		fmt.Printf("unable to spin up test server: %v", err)
	}

	var c net.Conn
	go func() {
		c, err = l.Accept()
		if err != nil {
			fmt.Printf("unable to accept connection on test server: %v", err)
		}
	}()
	defer l.Close()

	if conn, err = net.Dial("tcp4", ":54321"); err != nil {
		t.Fatalf("unable to connect to test port: %v\n", err)
	}
	defer conn.Close()

	// Seems to be a race condition from when the connection is made to when we try to
	// read from the bufio.Reader, occasionally causing nil pointer. Minor sleep seems to avoid that
	time.Sleep(time.Millisecond * 50)

	client := NewTestClientState(c, stateAccountMenu)
	client.AccountInfo.Account = acct.AccountName()
	fakePlayer, err := LoadPlayer(client, p)
	if err != nil {
		t.Fatalf("failed loading test player: %v\n", err)
	}
	client.AccountInfo.Player = fakePlayer.GetName()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		err = callAccountMenu(client, &wg)
	}()

	i := 0
	for i < 3 {
		select {
		case o := <-client.OutputStream:
			switch i {
			case 0:
				// Menu
				if _, e := conn.Write([]byte("L\n")); e != nil {
					t.Fatalf("unable to write to connection: %v\n", e)
				}
				break
			case 1:
				// character list
				assert.Equal(t, "<Y>TestPlayer</Y>: Level 1 TestClass", o)
				if _, e := conn.Write([]byte("q\n")); e != nil {
					t.Fatalf("unable to write to connection: %v\n", e)
				}
				break
			case 2:
				// quit
				assert.Equal(t, "Disconnected!", o)
				break
			}
			i++
		}
	}
	// Wait for go routine to complete before assertions
	wg.Wait()

	assert.Nil(t, err)
}
