package tests

type MockDirection struct {
	name        string
	destination string
}

func NewMockDirection(dirName string, dest string) *MockDirection {
	return &MockDirection{
		name:        dirName,
		destination: dest,
	}
}

func (m *MockDirection) Name() string {
	return m.name
}

func (m *MockDirection) SetName(name string) {
	m.name = name
}

func (m *MockDirection) Destination() string {
	return m.destination
}

func (m *MockDirection) SetDestination(roomId string) {
	m.destination = roomId
}
