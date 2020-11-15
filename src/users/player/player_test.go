package player

import (
	"races"
	"testing"
	. "tests"

	"github.com/stretchr/testify/assert"
)

func NewTestPlayer() *Player {
	races.Races = make(map[string]*races.Race)

	client := NewMockClient()
	r := NewMockRace()
	races.Races["TestClass"] = r
	return &Player{
		Name:        "Mrbagginz",
		DisplayName: "Mrbagginz",
		Level:       1,
		Health:      100,
		Fatigue:     100,
		Power:       100,
		Buffs:       nil,
		Debuffs:     nil,
		Location:    "123",
		Client:      client,
		Account:     "testAccount",
		Title:       "Goblin Maester",
		RealmTitle:  "Goblin Jedi Master",
		Race:        r,
		PlayerStats: nil,
		Stance:      0,
		Skills:      nil,
		Spells:      nil,
	}
}

// These tests copy mob_test but will keep them to ensure embedded structs are working as expected
func TestPlayer_AdjustFatigue(t *testing.T) {
	p := NewTestPlayer()

	p.AdjustFatigue(-5)
	assert.Equal(t, p.GetFatigue(), int16(95))

	p.AdjustFatigue(5)
	assert.Equal(t, p.GetFatigue(), int16(100))
}

func TestPlayer_AdjustHealth(t *testing.T) {
	p := NewTestPlayer()

	p.AdjustHealth(-5)
	assert.Equal(t, p.GetHealth(), int16(95))

	p.AdjustHealth(5)
	assert.Equal(t, p.GetHealth(), int16(100))
}

func TestPlayer_AdjustPower(t *testing.T) {
	p := NewTestPlayer()

	p.AdjustPower(-5)
	assert.Equal(t, p.GetPower(), int16(95))

	p.AdjustPower(5)
	assert.Equal(t, p.GetPower(), int16(100))
}

func TestPlayer_SetDisplayName(t *testing.T) {
	p := NewTestPlayer()

	p.SetDisplayName("Schmitty McWibberManJensen")
	assert.Equal(t, p.GetDisplayName(), "Schmitty McWibberManJensen")
}

func TestPlayer_SetRace(t *testing.T) {
	p := NewTestPlayer()

	err := p.SetRace("TestClass")
	assert.Nil(t, err)
	assert.Equal(t, p.RaceName(), "TestClass")
}

func TestPlayer_SetRace_Failed(t *testing.T) {
	p := NewTestPlayer()

	err := p.SetRace("Orc")
	assert.NotNil(t, err)
}
