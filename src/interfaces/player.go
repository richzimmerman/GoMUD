package interfaces

type PlayerInterface interface {
	GetTitle() string
	SetTitle(s string)
	GetRealmTitle() string
	SetRealmTitle(s string)
	RaceName() string
	SetRace(s string) error
	Realm() int8
	AccountName() string
	SetAccountName(s string)
	GetLocation() string
	SetLocation(roomId string)
	GetName() string
	SetName(s string)
	GetDisplayName() string
	SetDisplayName(s string)
	GetLevel() int8
	AdjustLevel(i int8)
	GetHealth() int16
	AdjustHealth(i int16)
	GetFatigue() int16
	AdjustFatigue(i int16)
	GetPower() int16
	AdjustPower(i int16)
	Send(msg string)
}
