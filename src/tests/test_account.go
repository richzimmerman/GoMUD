package tests

import . "interfaces"

type MockAccount struct {
	name     string
	lastIp   string
	email    string
	chars    []string
	password string
	status   bool
}

func NewMockAccount() AccountInterface {
	return &MockAccount{
		name:     "TestAccount",
		lastIp:   "127.0.0.1",
		email:    "foo@test.com",
		chars:    []string{"TestPlayer"},
		password: "TestPassword",
		status:   false,
	}
}

func (m *MockAccount) AccountName() string {
	return m.name
}

func (m *MockAccount) LastIPLogged() string {
	return m.lastIp
}

func (m *MockAccount) EmailAddress() string {
	return m.email
}

func (m *MockAccount) GetCharacters() []string {
	return m.chars
}

func (m *MockAccount) GetPassword() string {
	return m.password
}

func (m *MockAccount) ChangePassword(s string) error {
	m.password = s
	return nil
}

func (m *MockAccount) LoggedInStatus() bool {
	return m.status
}

func (m *MockAccount) SetLoggedInStatus(b bool) {
	m.status = b
}
