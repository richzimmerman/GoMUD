package client

import (
	"bufio"
	"db"
	"fmt"
	"io"
	"net"
	"sync"
	"telnet"
	"users/account"

	"github.com/DATA-DOG/go-sqlmock"
)

type InputStreamInterface interface {
	io.Reader
}

type MockInputStream struct {
	bufio.Reader
}

func (m *MockInputStream) Read(b []byte) (int, error) {
	return len(b), nil
}

func NewMockInputReader() *MockInputStream {
	return &MockInputStream{}
}

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
		Name:     "Foo",
		Password: "TestPassword",
		LastIP:   "127.0.0.1",
		Email:    "foo@test.com",
	}
	return &account.Account{
		DBAccount:  dba,
		Characters: make([]string, 0),
	}
}

func NewTestClientState(conn net.Conn, state int8) *Client {
	io := bufio.NewReader(conn)
	return &Client{
		loggedIn:     false,
		state:        state,
		Connection:   nil,
		Telnet:       telnet.NewTelnet(conn, io),
		Name:         "",
		OutputStream: make(chan string),
		outMutex:     sync.Mutex{},
		In:           nil,
		AccountInfo:  &loggedInfo{Account: "TestAccount", Player: "TestPlayer"},
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
