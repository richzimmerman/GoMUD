package client

import (
	"db"
	"fmt"
	"lib/accounts"
	"net"
	"sync"
	"testing"
	. "tests"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func callLogin(client *Client, wg *sync.WaitGroup) (string, error) {
	defer wg.Done()
	return client.logIn()
}

func TestLogin_Success(t *testing.T) {
	var err error
	var acctName string
	var conn net.Conn

	mock, err := InitMockDB()
	if err != nil {
		t.Fatalf("error mocking database: %v", err)
	}
	defer db.DatabaseConnection.Connection.Close()

	acct := NewMockAccount()
	accounts.AddAccount(acct)
	defer accounts.RemoveAccount(acct.AccountName())

	accountNameQuery := sqlmock.NewRows([]string{"Name"}).AddRow("TestAccount")
	passQuery := sqlmock.NewRows([]string{"Password"}).AddRow("TestPassword")
	accountQuery := sqlmock.NewRows([]string{"Name", "Password", "LastIP", "Email"}).
		AddRow("TestAccount", "TestPassword", "127.0.0.1", "foo@test.com")
	characterQuery := sqlmock.NewRows([]string{"Name", "Account", "DisplayName", "Level", "Health", "Fatigue", "Power",
		"Title", "RealmTitle", "Race", "Stats", "Stance", "Skills", "Spells", "Buffs", "Debuffs"})

	mock.ExpectPrepare("SELECT").ExpectQuery().WithArgs("TestAccount").
		WillReturnRows(accountNameQuery)
	mock.ExpectQuery("SELECT").WillReturnRows(passQuery)
	mock.ExpectPrepare("SELECT \\* FROM Accounts").ExpectQuery().WithArgs().
		WillReturnRows(accountQuery)
	mock.ExpectQuery("SELECT \\* FROM Characters").WillReturnRows(characterQuery)

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

	client := NewTestClientState(c, stateLogin)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		acctName, err = callLogin(client, &wg)
	}()

	i := 0
	for i < 2 {
		select {
		case _ = <-client.OutputStream:
			switch i {
			case 0:
				// account
				if _, e := conn.Write([]byte("TestAccount\n")); e != nil {
					t.Fatalf("unable to write to connection: %v\n", e)
				}
				break
			case 1:
				// password
				if _, e := conn.Write([]byte("TestPassword\n")); e != nil {
					t.Fatalf("unable to write to connection: %v\n", e)
				}
				break
			}
			i++
		}
	}
	// Wait for go routine to complete before assertions
	wg.Wait()

	assert.Equal(t, "TestAccount", acctName)
	assert.Nil(t, err)
}
