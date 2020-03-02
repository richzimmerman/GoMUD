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

func callChangePassword(client *Client, wg *sync.WaitGroup) error {
	defer wg.Done()
	err := client.changePassword()
	return err
}

func TestChange_Password_Success(t *testing.T) {
	var conn net.Conn

	mock, err := InitMockDB()
	if err != nil {
		t.Fatalf("error mocking database: %v", err)
	}
	defer db.DatabaseConnection.Connection.Close()

	mock.ExpectPrepare("UPDATE Accounts").WillBeClosed().ExpectExec().
		WithArgs("newtestpassword", "Foo").WillReturnResult(sqlmock.NewResult(1, 1))

	client := NewTestClientState(stateAccountChangePassword)
	client.Account = NewMockAccount()

	l := NewTestServer(client)
	defer l.Close()

	if conn, err = net.Dial("tcp4", ":54321"); err != nil {
		t.Fatalf("unable to connect to test port: %v\n", err)
	}
	defer conn.Close()

	// Seems to be a race condition from when the connection is made to when callChangePassword tries to
	// read from the bufio Reader, occasionally causing nil pointer. Minor sleep seems to avoid that
	time.Sleep(time.Millisecond * 50)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		err = callChangePassword(client, &wg)
	}()

	// Prompt will ask for new password and confirmation via client.OutputStream so write new password on prompts.
	i := 0
	for i < 2 {
		select {
		case _ = <- client.OutputSteam:
			if _, e := conn.Write([]byte("newtestpassword\n")); e != nil {
				t.Fatalf("unable to write to connection: %v\n", e)
			}
			i++
		}
	}
	// Wait for go routine to complete before asserting new password value
	wg.Wait()

	assert.Equal(t, "newtestpassword", client.Account.Password)
	assert.Nil(t, err)
}

func TestChange_Password_Mismatch(t *testing.T) {
	client := NewTestClientState(stateAccountChangePassword)
	client.Account = NewMockAccount()

	l := NewTestServer(client)
	defer l.Close()

	var err error
	var conn net.Conn
	if conn, err = net.Dial("tcp4", ":54321"); err != nil {
		t.Fatalf("unable to connect to test port: %v\n", err)
	}
	defer conn.Close()

	// Seems to be a race condition from when the connection is made to when callChangePassword tries to
	// read from the bufio Reader, occasionally causing nil pointer. Minor sleep seems to avoid that
	time.Sleep(time.Millisecond * 50)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		err = callChangePassword(client, &wg)
	}()

	// Prompt will ask for new password and confirmation via client.OutputStream so write new password on prompts.
	i := 0
	for i < 2 {
		select {
		case _ = <- client.OutputSteam:
			if i == 0 {
				if _, e := conn.Write([]byte("newtestpassword\n")); e != nil {
					t.Fatalf("unable to write to connection: %v\n", e)
				}
			} else {
				if _, e := conn.Write([]byte("notnewtestpassword\n")); e != nil {
					t.Fatalf("unable to write to connection: %v\n", e)
				}
			}
			i++
		}
	}
	// Wait for go routine to complete before asserting new password value
	wg.Wait()

	assert.Equal(t, "TestPassword", client.Account.Password)
	assert.NotNil(t, err)
	assert.Equal(t, "password change failed, passwords do not match", err.Error())
}