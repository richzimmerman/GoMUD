package mobs

import (
	"db"
	"encoding/json"
	"fmt"
	"races"
	"strings"
	"utils"
	"world/realms"
)

const (
	stanceParry      = 0
	stanceDefensive  = 1
	stanceNormal     = 2
	stanceAggressive = 3
)

type Player struct {
	account    string
	title      string
	realmTitle string
	race       *races.Race
	stats      map[string]int8
	stance     int8
	skills     map[string]*interface{}
	spells     map[string]*interface{}
	*Mob
}

func NewPlayer(account string, name string, race string, realm int8) (*Player, error) {
	r, ok := races.Races[race]
	if !ok {
		return nil, utils.Error(fmt.Errorf("invalid race, or race not loaded: %s", race))
	}
	if r.Realm != realm {
		return nil, utils.Error(fmt.Errorf("realm mismatch, %s is not part of realm %s",
			race, realms.Realms[realm]))
	}
	mob := &Mob{
		name:        name,
		displayName: name,
		level:       1,
		health:      r.DefaultHealth,
		fatigue:     r.DefaultFatigue,
		power:       r.DefaultPower,
		buffs:       nil,
		debuffs:     nil,
		location:    r.StartingRoom,
	}
	p := &Player{
		account:    account,
		title:      r.DefaultTitle,
		realmTitle: "",
		race:       r,
		stats:      r.DefaultStats,
		stance:     0,
		skills:     nil,
		spells:     nil,
		Mob:        mob,
	}
	err := SyncPlayer(p)
	if err != nil {
		return nil, utils.Error(err)
	}
	return p, nil
}

func SyncPlayer(p *Player) error {
	/* Converts a *Player to a *DBPlayer and saves it in the DB. */
	var pStats string
	s, err := json.Marshal(p.stats)
	if err != nil {
		return utils.Error(err)
	}
	pStats = string(s)

	sk, err := utils.JsonMarshalMapToArray(p.skills)
	if err != nil {
		return utils.Error(err)
	}

	sp, err := utils.JsonMarshalMapToArray(p.spells)
	if err != nil {
		return utils.Error(err)
	}

	bf, err := json.Marshal(p.buffs)
	if err != nil {
		return utils.Error(err)
	}
	buffs := string(bf)
	if buffs == "null" {
		buffs = "{}"
	}

	dbf, err := json.Marshal(p.debuffs)
	if err != nil {
		return utils.Error(err)
	}
	debuffs := string(dbf)
	if debuffs == "null" {
		debuffs = "{}"
	}

	d := &db.DBPlayer{
		Name:        p.Name(),
		Account:     p.Account(),
		DisplayName: p.DisplayName(),
		Level:       p.Level(),
		Health:      p.Health(),
		Fatigue:     p.Fatigue(),
		Power:       p.Power(),
		Title:       p.Title(),
		RealmTitle:  p.RealmTitle(),
		Race:        p.Race(),
		Stats:       pStats,
		Stance:      0,
		Skills:      sk,
		Spells:      sp,
		Buffs:       buffs,
		Debuffs:     debuffs,
		Location:    p.location,
	}

	err = db.DatabaseConnection.SavePlayer(d)
	return utils.Error(err)
}

func LoadPlayer(p *db.DBPlayer) (*Player, error) {
	// TODO: this for buffs and debuffs. Remember to parse JSON first.
	//b := make(map[string]*spells.Buff)
	//for _, buff := range p.Buffs {
	//	b[buff], err = spells.LoadBuff(buff, 0)
	//}
	stats := make(map[string]int8)
	err := json.Unmarshal([]byte(p.Stats), &stats)
	if err != nil {
		return nil, utils.Error(err)
	}

	buffs := make(map[string]string)
	err = json.Unmarshal([]byte(p.Buffs), &buffs)
	if err != nil {
		return nil, utils.Error(err)
	}

	debuffs := make(map[string]string)
	err = json.Unmarshal([]byte(p.Buffs), &debuffs)
	if err != nil {
		return nil, utils.Error(err)
	}

	mob := &Mob{
		p.Name,
		p.DisplayName,
		p.Level,
		p.Health,
		p.Fatigue,
		p.Power,
		buffs,
		debuffs,
		"",
	}
	// TODO: parse JSON for stats, skills, spells
	race, ok := races.Races[p.Race]
	if !ok {
		return nil, utils.Error(fmt.Errorf("invalid race, or race not loaded: %s", p.Race))
	}
	return &Player{
		p.Account,
		p.Title,
		p.RealmTitle,
		race,
		stats,
		0,
		nil,
		nil,
		mob,
	}, nil
}

func (p *Player) Title() string {
	return p.title
}

func (p *Player) SetTitle(title string) {
	p.title = title
}

func (p *Player) RealmTitle() string {
	return p.realmTitle
}

func (p *Player) RealmSetTitle(realmTitle string) {
	p.realmTitle = realmTitle
}

func (p *Player) Race() string {
	return p.race.Name
}

func (p *Player) SetRace(race string) error {
	// TODO: This one seems like it could cause problems w/ skills and such. Will require testing once we're in game and doing stuff.
	race = strings.Title(race)
	r, ok := races.Races[race]
	if !ok {
		return utils.Error(fmt.Errorf("unable to set race: %s", race))
	}
	p.race = r
	return nil
}

func (p *Player) Account() string {
	return p.account
}

func (p *Player) SetAccount(account string) error {
	p.account = account
	err := SyncPlayer(p)
	return utils.Error(err)
}
