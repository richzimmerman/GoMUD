package tests

type MockClient struct {
	accountName string
	player      string
}

func (m *MockClient) AssociatedAccount() (string, error) {
	return m.accountName, nil
}

func (m *MockClient) SetAssociatedAccount(name string) {
	m.accountName = name
}

func (m *MockClient) AssociatedPlayer() string {
	return m.player
}

func (m *MockClient) SetAssociatedPlayer(name string) {
	m.player = name
}

func (m *MockClient) GameLoop() error {
	return nil
}

func (m *MockClient) Out(msg string) {
	// do nothing
}

func (m *MockClient) GetRemoteAddress() string {
	return "127.0.0.1:12345"
}

func NewMockClient() *MockClient {
	return &MockClient{
		accountName: "testAccount",
		player:      "testPlayer",
	}
}
