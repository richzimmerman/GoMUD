package client

import (
	"db"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"net"
	"sync"
	"testing"
	"time"
)

func callCreateAccount(client *Client, wg *sync.WaitGroup) (bool, error) {
	defer wg.Done()
	return client.createAccountPrompt()
}

func TestAccountCreation_Success(t *testing.T) {
	var conn net.Conn

	expectedAccountName := "Fakeaccount"
	expectedPassword := "testPassword"
	expectedLastIP := "127.0.0.1"
	expectedEmail := "foo@test.com"

	mock, err := InitMockDB()
	if err != nil {
		t.Fatalf("error mocking database: %v", err)
	}
	defer db.DatabaseConnection.Connection.Close()

	mock.ExpectPrepare("INSERT INTO Accounts").WillBeClosed().ExpectExec().
		WithArgs(expectedAccountName, expectedPassword, expectedLastIP, expectedEmail).
		WillReturnResult(sqlmock.NewResult(1, 1))

	client := NewTestClientState(stateAccountChangePassword)

	l := NewTestServer(client)
	defer l.Close()

	if conn, err = net.Dial("tcp4", ":54321"); err != nil {
		t.Fatalf("unable to connect to test port: %v\n", err)
	}
	defer conn.Close()

	// Seems to be a race condition from when the connection is made to when callChangePassword tries to
	// read from the bufio Reader, occasionally causing nil pointer. Minor sleep seems to avoid that
	time.Sleep(time.Millisecond * 50)

	var ok bool
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		ok, err = callCreateAccount(client, &wg)
	}()

	// For sanity's sake, make sure account is nil before creating one. client.Account struct should be created along
	// with the account in the database.
	assert.Nil(t, client.Account)

	// Prompt will ask for new password and confirmation via client.OutputStream so write new password on prompts.
	i := 0
	for i < 5 {
		select {
		case _ = <- client.OutputSteam:
			switch i {
			case 0:
				//account
				if _, e := conn.Write([]byte(expectedAccountName + "\n")); e != nil {
					t.Fatalf("unable to write to connection: %v\n", e)
				}
				break
			case 1:
				// account confirm
				if _, e := conn.Write([]byte("Y\n")); e != nil {
					t.Fatalf("unable to write to connection: %v\n", e)
				}
				break
			case 2:
				//password
				if _, e := conn.Write([]byte(expectedPassword + "\n")); e != nil {
					t.Fatalf("unable to write to connection: %v\n", e)
				}
				break
			case 3:
				// password confirm
				if _, e := conn.Write([]byte(expectedPassword + "\n")); e != nil {
					t.Fatalf("unable to write to connection: %v\n", e)
				}
				break
			case 4:
				// email address
				if _, e := conn.Write([]byte(expectedEmail + "\n")); e != nil {
					t.Fatalf("unable to write to connection: %v\n", e)
				}
				break
			}
			i++
		}
	}
	// Wait for go routine to complete before assertions
	wg.Wait()

	assert.True(t, ok)
	assert.Nil(t, err)
	assert.Equal(t, expectedAccountName, client.Account.Name)
	assert.Equal(t, expectedPassword, client.Account.Password)
	assert.Equal(t, expectedLastIP, client.Account.LastIP)
	assert.Equal(t, expectedEmail, client.Account.Email)
}
