package spells

import "time"

type Debuff struct {
	Name     string
	Duration time.Duration
}

// TODO: interface
func LoadDebuff(name string, duration time.Duration) (*Debuff, error) {
	// TODO: Load buff/debuff from database
	return &Debuff{
		Name:     name,
		Duration: duration,
	}, nil
}
