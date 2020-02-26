package spells

import "time"

type Buff struct {
	Name     string
	Duration time.Duration
}

// TODO: interface
func LoadBuff(name string, duration time.Duration) (*Buff, error) {
	// TODO: Load buff/debuff from database
	return &Buff{
		Name:     name,
		Duration: duration,
	}, nil
}
