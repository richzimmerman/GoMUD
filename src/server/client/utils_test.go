package client

import (
	"account"
	"bufio"
	"db"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"mobs"
	"net"
	"sync"
)

func InitMockDB() (sqlmock.Sqlmock, error) {
	d, mock, err := sqlmock.New()
	if err != nil {
		return nil, err
	}

	db.DatabaseConnection = &db.DbConnection{
		Connection: d,
	}

	return mock, nil
}

func NewMockAccount() *account.Account {
	dba := &db.DBAccount{
		Name: "Foo",
		Password: "TestPassword",
		LastIP: "127.0.0.1",
		Email: "foo@test.com",
	}
	return &account.Account{
		DBAccount:  dba,
		Characters: make(map[string]*mobs.Player),
	}
}

func NewTestClientState(state int8) *Client {
	return &Client{
		loggedIn:    false,
		state:       state,
		Connection:  nil,
		Telnet:      nil,
		Name:        "",
		Player:      nil,
		Account:     nil,
		OutputSteam: make(chan string),
		outMutex:    sync.Mutex{},
		In:          nil,
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
