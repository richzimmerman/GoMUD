package client

import (
	"bufio"
	"db"
	"io"
	"net"
	"sync"
	"telnet"

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

func NewTestClientState(conn net.Conn, state int8) *Client {
	in := bufio.NewReader(conn)
	return &Client{
		loggedIn:     false,
		state:        state,
		Connection:   conn,
		Telnet:       telnet.NewTelnet(conn, in),
		Name:         "",
		OutputStream: make(chan string),
		outMutex:     sync.Mutex{},
		In:           in,
		AccountInfo:  &loggedInfo{Account: "TestAccount", Player: "TestPlayer"},
	}
}
