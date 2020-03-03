package db

type DBAccount struct {
	Name     string
	Password string
	LastIP   string
	Email    string
}

type DBPlayer struct {
	Name        string
	Account     string
	DisplayName string
	Level       int8
	Health      int16
	Fatigue     int16
	Power       int16
	Title       string
	RealmTitle  string
	Race        string
	Stats       string // JSON String of stats
	Stance      int8
	Skills      string // JSON Array: List of skills to load later ["skill1", "skill2", "skill3"]
	Spells      string // Same as above
	Buffs       string // JSON array of objects [{"name": "buff1", "duration":12345}, {}, {}]
	Debuffs     string // Same as above
	Location    string // Room ID
}

type DBRace struct {
	Name           string
	Realm          int8
	Type           int8
	SkillList      string // JSON array
	Description    string
	DefaultHealth  int16
	DefaultFatigue int16
	DefaultPower   int16
	StartingRoom   string // room ID
	DefaultTitle   string
	DefaultStats   string // JSON String of stats
}
