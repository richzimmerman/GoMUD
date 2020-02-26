package mobs

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func NewTestPlayer() *Player {
	mob := &Mob{
		"Mrbagginz",
		"Mrbagginz",
		100,
		100,
		100,
		nil,
		nil,
	}
	return &Player{
		"Goblin Maester",
		"Goblin Jedi Master",
		"Goblin",
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

func TestPlayer_SetClass(t *testing.T) {
	p := NewTestPlayer()

	p.SetClass("Orc")
	assert.Equal(t, p.Class(), "Orc")
}
