package client

import (
	"bufio"
	"db"
	"fmt"
	. "lib/players"
	"net"
	"races"
	"sync"
	"testing"
	"tests"
	"time"

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

	mockClient := tests.NewMockClient()
	fakePlayer, err := LoadPlayer(mockClient, p)
	if err != nil {
		t.Fatalf("failed loading test player: %v\n", err)
	}

	l, err := net.Listen("tcp4", ":54321")
	if err != nil {
		fmt.Printf("unable to spin up test server: %v", err)
	}

	go func() {
		_, err := l.Accept()
		if err != nil {
			fmt.Printf("unable to accept connection on test server: %v", err)
		}
	}()
	defer l.Close()

	if conn, err = net.Dial("tcp4", ":54321"); err != nil {
		t.Fatalf("unable to connect to test port: %v\n", err)
	}
	defer conn.Close()

	client := NewTestClientState(conn, stateAccountMenu)
	client.Connection = conn
	client.In = bufio.NewReader(conn)
	client.AccountInfo.Account = NewMockAccount().AccountName()
	client.AccountInfo.Player = fakePlayer.GetName()

	// Seems to be a race condition from when the connection is made to when callChangePassword tries to
	// read from the bufio Reader, occasionally causing nil pointer. Minor sleep seems to avoid that
	time.Sleep(time.Millisecond * 50)

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
				// q
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
