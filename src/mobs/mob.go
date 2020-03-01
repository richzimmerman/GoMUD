package mobs

import "spells"

/*
This is the base race Mob for any "living" being in the game, ie. Players, NPCs, and creatures
*/
type Mob struct {
	name        string
	displayName string
	level       int8
	health      int16
	fatigue     int16
	power       int16
	// TODO: might wants buffs to be a struct
	buffs map[string]*spells.Buff
	// TODO: might want debuffs to be a struct
	debuffs map[string]*spells.Debuff
}

func (m *Mob) Name() string {
	return m.name
}

func (m *Mob) DisplayName() string {
	return m.displayName
}

func (m *Mob) SetDisplayName(name string) {
	m.displayName = name
}

func (m *Mob) Level() int8 {
	return m.level
}

func (m *Mob) AdjustLevel(i int8) {
	m.level += i
}

func (m *Mob) Health() int16 {
	return m.health
}

func (m *Mob) AdjustHealth(i int16) {
	m.health += i
}

func (m *Mob) Fatigue() int16 {
	return m.fatigue
}

func (m *Mob) AdjustFatigue(i int16) {
	m.fatigue += i
}

func (m *Mob) Power() int16 {
	return m.power
}

func (m *Mob) AdjustPower(i int16) {
	m.power += i
}
