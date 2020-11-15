package player

import (
	"fmt"
	. "interfaces"
	"logger"
	"races"
	"strings"
	"utils"
)

var log = logger.NewLogger()

const (
	stanceParry      = 0
	stanceDefensive  = 1
	stanceNormal     = 2
	stanceAggressive = 3
)

// Player struct
type Player struct {
	Name        string
	DisplayName string
	Level       int8
	Health      int16
	Fatigue     int16
	Power       int16
	Buffs       map[string]string
	Debuffs     map[string]string
	Location    string
	Client      ClientInterface
	Account     string
	Title       string
	RealmTitle  string
	Race        *races.Race
	PlayerStats map[string]int8
	Stance      int8
	Skills      map[string]*interface{}
	Spells      map[string]*interface{}
}

func (p *Player) GetName() string {
	return p.Name
}

func (p *Player) SetName(name string) {
	p.Name = name
}

func (p *Player) GetDisplayName() string {
	return p.DisplayName
}

func (p *Player) SetDisplayName(name string) {
	p.DisplayName = name
}

func (p *Player) GetLevel() int8 {
	return p.Level
}

func (p *Player) AdjustLevel(i int8) {
	lvl := p.Level + 1
	if lvl < 0 {
		p.Level = 0
	} else {
		p.Level = lvl
	}
}

func (p *Player) GetHealth() int16 {
	return p.Health
}

func (p *Player) AdjustHealth(i int16) {
	h := p.Health + i
	if h < 0 {
		p.Health = 0
	} else {
		p.Health = h
	}
}

func (p *Player) GetFatigue() int16 {
	return p.Fatigue
}

func (p *Player) AdjustFatigue(i int16) {
	fat := p.Fatigue + i
	if fat < 0 {
		p.Fatigue = 0
	} else {
		p.Fatigue = fat
	}
}

func (p *Player) GetPower() int16 {
	return p.Power
}

func (p *Player) AdjustPower(i int16) {
	pwr := p.Power + i
	if pwr < 0 {
		p.Power = 0
	} else {
		p.Power = pwr
	}
}

func (p *Player) GetLocation() string {
	return p.Location
}

func (p *Player) SetLocation(roomId string) {
	p.Location = roomId
}

// Title getter
func (p *Player) GetTitle() string {
	return p.Title
}

// SetTitle sets player's title
func (p *Player) SetTitle(title string) {
	p.Title = title
}

// RealmTitle getter
func (p *Player) GetRealmTitle() string {
	return p.RealmTitle
}

// SetRealmTitle sets player's PVP realm title
func (p *Player) SetRealmTitle(realmTitle string) {
	p.RealmTitle = realmTitle
}

// Race getter
func (p *Player) RaceName() string {
	return p.Race.Name
}

// SetRace sets the players races
func (p *Player) SetRace(race string) error {
	// TODO: This one seems like it could cause problems w/ skills and such. Will require testing once we're in game and doing stuff.
	race = strings.Title(race)
	r, ok := races.Races[race]
	if !ok {
		return utils.Error(fmt.Errorf("unable to set race: %s", race))
	}
	p.Race = r
	return nil
}

// Realm is the integer representation of the realm a player belongs to
func (p *Player) Realm() int8 {
	return p.Race.Realm
}

// Account returns name of the account player belongs to
func (p *Player) AccountName() string {
	return p.Account
}

func (p *Player) SetAccountName(name string) {
	p.Account = name
}

// Send sends output to the user, passed through the output formatter for ANSI coloring
func (p *Player) Send(msg string) {
	p.Client.Out(msg)
}
