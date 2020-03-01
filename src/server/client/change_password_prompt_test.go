package client

import (
	"account"
	"bufio"
	"db"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"net"
	"sync"
	"testing"
	"time"
)

func NewMockAccount() *account.Account {
	dba := &db.DBAccount{
		Name: "Foo",
		Password: "TestPassword",
		LastIP: "127.0.0.1",
		Email: "foo@test.com",
	}
	return &account.Account{
		DBAccount:  dba,
		Characters: nil,
	}
}

func NewTestServer(client *Client) net.Listener {
	var err error
	l, err := net.Listen("tcp4", ":54321")
	if err != nil {
		fmt.Printf("unable to spin up test server: %v", err)
	}

	go func() {
		c, err := l.Accept()
		if err != nil {
			fmt.Printf("unable to accept connection on test server: %v", err)
		}
		client.Connection = c
		client.In = bufio.NewReader(c)
	}()
	return l
}

func callChangePassword(client *Client, wg *sync.WaitGroup) error {
	defer wg.Done()
	err := client.changePassword()
	return err
}

func TestChange_Password_Success(t *testing.T) {
	d, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error mocking database: %v", err)
	}
	defer d.Close()

	db.DatabaseConnection = &db.DbConnection{
		Connection: d,
	}

	mock.ExpectPrepare("UPDATE Accounts").WillBeClosed().ExpectExec().
		WithArgs("newtestpassword", "Foo").WillReturnResult(sqlmock.NewResult(1, 1))

	client := &Client{
		loggedIn:    false,
		state:       stateAccountChangePassword,
		Connection:  nil,
		Telnet:      nil,
		Name:        "",
		Player:      nil,
		Account:     nil,
		OutputSteam: make(chan string),
		outMutex:    sync.Mutex{},
		In:          nil,
	}
	client.Account = NewMockAccount()

	l := NewTestServer(client)
	defer l.Close()

	var conn net.Conn
	if conn, err = net.Dial("tcp4", ":54321"); err != nil {
		t.Fatalf("unable to connect to test port: %v\n", err)
	}
	defer conn.Close()

	// Seems to be a race condition from when the connection is made to when callChangePassword tries to
	// read from the bufio Reader, occasionally causing nil pointer. Minor sleep seems to avoid that
	time.Sleep(time.Millisecond * 250)

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
	client := &Client{
		loggedIn:    false,
		state:       stateAccountChangePassword,
		Connection:  nil,
		Telnet:      nil,
		Name:        "",
		Player:      nil,
		Account:     nil,
		OutputSteam: make(chan string),
		outMutex:    sync.Mutex{},
		In:          nil,
	}
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
	time.Sleep(time.Millisecond * 250)

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