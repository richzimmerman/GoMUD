package mobs

import (
	"crypto/rand"
	"fmt"
)

// TODO: onDestroy() Remove object from rooms and stuff

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
	buffs       map[string]string
	debuffs     map[string]string
	location    string
	guid        string
	// TODO: Alliance (realm) to indicate friendly and enemy NPCs or Players
}

func (m *Mob) GetGUID() string {
	return m.guid
}

func (m *Mob) GenerateGUID() error {
	if m.guid != "" {
		return fmt.Errorf("cannot modify an existing mobs GUID")
	}
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return err
	}
	m.guid = fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return nil
}

func (m *Mob) String() string {
	return m.DisplayName()
}

func (m *Mob) Name() string {
	return m.name
}

func (m *Mob) SetName(name string) {
	m.name = name
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

func (m *Mob) Location() string {
	return m.location
}

func (m *Mob) SetLocation(roomId string) {
	m.location = roomId
}
