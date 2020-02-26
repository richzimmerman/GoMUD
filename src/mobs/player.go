package mobs

import (
	"db"
	"skills"
	"spells"
)

const (
	stanceParry      = 0
	stanceDefensive  = 1
	stanceNormal     = 2
	stanceAggressive = 3
)

type Player struct {
	title      string
	realmTitle string
	race       string
	// TODO: might want stats to be a struct?
	stats map[string]int8
	stance  int8
	skills  map[string]*skills.Skill
	spells  map[string]*spells.Spell
	*Mob
}

// TODO: figure out how to not have circular imports with db
func LoadPlayer(p *db.DBCharacter) (*Player, error) {
	// TODO: Load player data from database with DBPlayer struct
	// TODO: this for buffs and debuffs. Remember to parse JSON first.
	//b := make(map[string]*spells.Buff)
	//for _, buff := range p.Buffs {
	//	b[buff], err = spells.LoadBuff(buff, 0)
	//}
	mob := &Mob{
		p.Name,
		p.DisplayName,
		p.Health,
		p.Fatigue,
		p.Power,
		nil,
		nil,
	}
	// TODO: parse JSON for stats, skills, spells
	return &Player{
		p.Title,
		p.RealmTitle,
		p.Race,
		nil,
		0,
		nil,
		nil,
		mob,
	}, nil
}

func (p *Player) Race() string {
	return p.race
}

func (p *Player) SetRace(class string) {
	p.race = class
}
