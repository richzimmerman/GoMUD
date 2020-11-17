package tests

import (
	. "interfaces"
	"time"
)

type MockSession struct {
	player    PlayerInterface
	client    ClientInterface
	lastInput int64
}

func NewMockSession(p PlayerInterface, c ClientInterface) SessionInterface {
	return &MockSession{
		player:    p,
		client:    c,
		lastInput: time.Now().Unix(),
	}
}

func (m *MockSession) Player() PlayerInterface {
	return m.player
}

func (m *MockSession) Client() ClientInterface {
	return m.client
}

func (m *MockSession) LastInput() int64 {
	return m.lastInput
}

func (m *MockSession) InputReceived(timestamp int64) {
	m.lastInput = timestamp
}
