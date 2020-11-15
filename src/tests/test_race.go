package tests

import "races"

func NewMockRace() *races.Race {
	return &races.Race{
		Name:           "TestClass",
		Realm:          0,
		Type:           0,
		SkillList:      nil,
		Description:    "",
		DefaultHealth:  0,
		DefaultFatigue: 0,
		DefaultPower:   0,
		StartingRoom:   "",
		DefaultTitle:   "",
		DefaultStats:   make(map[string]int8),
	}
}
