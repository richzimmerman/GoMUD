package mobs

import (
	"github.com/stretchr/testify/assert"
	"races"
	"testing"
)

func NewTestPlayer() *Player {
	races.Races = make(map[string]*races.Race)

	r := &races.Race{
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
	races.Races["TestClass"] = r

	mob := &Mob{
		"Mrbagginz",
		"Mrbagginz",
		1,
		100,
		100,
		100,
		nil,
		nil,
		"",
	}
	return &Player{
		"Goblin Maester",
		"Goblin Jedi Master",
		"Goblin",
		nil,
		nil,
		0,
		nil,
		nil,
		mob,
	}
}

// These tests copy mob_test but will keep them to ensure embedded structs are working as expected
func TestPlayer_AdjustFatigue(t *testing.T) {
	p := NewTestPlayer()

	p.AdjustFatigue(-5)
	assert.Equal(t, p.Fatigue(), int16(95))

	p.AdjustFatigue(5)
	assert.Equal(t, p.Fatigue(), int16(100))
}

func TestPlayer_AdjustHealth(t *testing.T) {
	p := NewTestPlayer()

	p.AdjustHealth(-5)
	assert.Equal(t, p.Health(), int16(95))

	p.AdjustHealth(5)
	assert.Equal(t, p.Health(), int16(100))
}

func TestPlayer_AdjustPower(t *testing.T) {
	p := NewTestPlayer()

	p.AdjustPower(-5)
	assert.Equal(t, p.Power(), int16(95))

	p.AdjustPower(5)
	assert.Equal(t, p.Power(), int16(100))
}

func TestPlayer_SetDisplayName(t *testing.T) {
	p := NewTestPlayer()

	p.SetDisplayName("Schmitty McWibberManJensen")
	assert.Equal(t, p.DisplayName(), "Schmitty McWibberManJensen")
}

func TestPlayer_SetRace(t *testing.T) {
	p := NewTestPlayer()

	err := p.SetRace("TestClass")
	assert.Nil(t, err)
	assert.Equal(t, p.Race(), "TestClass")
}

func TestPlayer_SetRace_Failed(t *testing.T) {
	p := NewTestPlayer()

	err := p.SetRace("Orc")
	assert.NotNil(t, err)
}