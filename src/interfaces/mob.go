package interfaces

type MobInterface interface {
	GetGUID() string
	GenerateGUID() error
	Name() string
	SetName(name string)
	DisplayName() string
	SetDisplayName(name string)
	Level() int8
	AdjustLevel(i int8)
	Health() int16
	AdjustHealth(i int16)
	Fatigue() int16
	AdjustFatigue(i int16)
	Power() int16
	AdjustPower(i int16)
	Location() string
	SetLocation(roomId string)
}
