package mobs

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func NewTestMob() *Mob {
	return &Mob{
		"Mrbagginz",
		"Mrbagginz",
		100,
		100,
		100,
		nil,
		nil,
	}
}

func TestMob_AdjustFatigue(t *testing.T) {
	p := NewTestMob()

	p.AdjustFatigue(-5)
	assert.Equal(t, p.Fatigue(), int16(95))

	p.AdjustFatigue(5)
	assert.Equal(t, p.Fatigue(), int16(100))
}

func TestMob_AdjustHealth(t *testing.T) {
	p := NewTestMob()

	p.AdjustHealth(-5)
	assert.Equal(t, p.Health(), int16(95))

	p.AdjustHealth(5)
	assert.Equal(t, p.Health(), int16(100))
}

func TestMob_AdjustPower(t *testing.T) {
	p := NewTestMob()

	p.AdjustPower(-5)
	assert.Equal(t, p.Power(), int16(95))

	p.AdjustPower(5)
	assert.Equal(t, p.Power(), int16(100))
}

func TestMob_SetDisplayName(t *testing.T) {
	p := NewTestMob()

	p.SetDisplayName("Schmitty McWibberManJensen")
	assert.Equal(t, p.DisplayName(), "Schmitty McWibberManJensen")
}