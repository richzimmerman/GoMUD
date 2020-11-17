package tests

type MockPlayer struct {
	title       string
	realmTitle  string
	race        string
	realm       int8
	account     string
	location    string
	name        string
	displayName string
	level       int8
	health      int16
	fatigue     int16
	power       int16
	fakeOutput  string
}

// NewMockPlayer returns a PlayerInterface that includes a fakeOutput string slice so when the Send function
//   is called, we can read and assert the expected output.
func NewMockPlayer() *MockPlayer {
	return &MockPlayer{
		title:       "TestTitle",
		realmTitle:  "Test Realm Title",
		race:        "Test Goblin",
		realm:       0,
		account:     "TestAccount",
		location:    "0",
		name:        "TestPlayer",
		displayName: "TestPlayer",
		level:       1,
		health:      100,
		fatigue:     100,
		power:       100,
		fakeOutput:  "",
	}
}

func (m *MockPlayer) GetTitle() string {
	return m.title
}

func (m *MockPlayer) SetTitle(s string) {
	m.title = s
}

func (m *MockPlayer) GetRealmTitle() string {
	return m.realmTitle
}

func (m *MockPlayer) SetRealmTitle(s string) {
	m.realmTitle = s
}

func (m *MockPlayer) RaceName() string {
	return m.race
}

func (m *MockPlayer) SetRace(s string) error {
	m.race = s
	return nil
}

func (m *MockPlayer) Realm() int8 {
	return m.realm
}

func (m *MockPlayer) AccountName() string {
	return m.account
}

func (m *MockPlayer) SetAccountName(s string) {
	m.account = s
}

func (m *MockPlayer) GetLocation() string {
	return m.location
}

func (m *MockPlayer) SetLocation(roomId string) {
	m.location = roomId
}

func (m *MockPlayer) GetName() string {
	return m.name
}

func (m *MockPlayer) SetName(s string) {
	m.name = s
}

func (m *MockPlayer) GetDisplayName() string {
	return m.displayName
}

func (m *MockPlayer) SetDisplayName(s string) {
	m.displayName = s
}

func (m *MockPlayer) GetLevel() int8 {
	return m.level
}

func (m *MockPlayer) AdjustLevel(i int8) {
	m.level += i
}

func (m *MockPlayer) GetHealth() int16 {
	return m.health
}

func (m *MockPlayer) AdjustHealth(i int16) {
	m.health += i
}

func (m *MockPlayer) GetFatigue() int16 {
	return m.fatigue
}

func (m *MockPlayer) AdjustFatigue(i int16) {
	m.fatigue += i
}

func (m *MockPlayer) GetPower() int16 {
	return m.power
}

func (m *MockPlayer) AdjustPower(i int16) {
	m.power += i
}

func (m *MockPlayer) Send(msg string) {
	m.fakeOutput = msg
}

func (m *MockPlayer) GetOutput() string {
	return m.fakeOutput
}
