package races

import (
	"db"
	"encoding/json"
	"fmt"
	"skills"
	"utils"
)

var Races map[string]*Race

type Race struct {
	Name           string
	Realm          int8
	Type           int8 // TODO: enum style int for race types (thief, crafty, spell caster, warrior, hybrid, tank) to effect skills/stats
	SkillList      map[string]*RaceSkill
	Description    string
	DefaultHealth  int16
	DefaultFatigue int16
	DefaultPower   int16
	StartingRoom   string // room ID
	DefaultTitle   string
	DefaultStats   map[string]int8
}

type RaceSkill struct {
	Skill          *skills.Skill
	LevelAvailable int8
}

func LoadRaces() error {
	fmt.Println("Loading Races.")
	Races = make(map[string]*Race)
	races, err := db.DatabaseConnection.LoadRaces()
	if err != nil {
		return utils.Error(err)
	}
	for _, race := range races {
		r, err := loadRace(race)
		if err != nil {
			return utils.Error(err)
		}
		Races[r.Name] = r
	}
	return nil
}

func loadRace(r *db.DBRace) (*Race, error) {
	// TODO: Skills
	// Unmarshal JSON string, iterate over skill list, load skill pointers into map
	// skills := ^that
	race := &Race{
		Name:        r.Name,
		Realm:       r.Realm,
		Type:        r.Type,
		SkillList:   nil,
		Description: r.Description,
		DefaultHealth: r.DefaultHealth,
		DefaultFatigue: r.DefaultFatigue,
		DefaultPower: r.DefaultPower,
		StartingRoom: r.StartingRoom,
		DefaultTitle: r.DefaultTitle,
		DefaultStats: make(map[string]int8),
	}

	err := json.Unmarshal([]byte(r.DefaultStats), &race.DefaultStats)
	if err != nil {
		return nil, utils.Error(err)
	}
	return race, nil
}

func NewRace(name string, realm int8, t int8) (*Race, error) {
	d := &db.DBRace{
		Name:        name,
		Realm:       realm,
		Type:        t,
		SkillList:   "[]",
		Description: "Default description.",
	}
	err := db.DatabaseConnection.CreateRace(d)
	if err != nil {
		return nil, utils.Error(err)
	}
	r, err := loadRace(d)
	if err != nil {
		return nil, utils.Error(err)
	}
	return r, nil
}

func (r *Race) AddSkill(skill *skills.Skill, unlockedAt int8) {
	r.SkillList[skill.Name] = &RaceSkill{
		Skill:          skill,
		LevelAvailable: unlockedAt,
	}
}

func (r *Race) RemoveSkill(skillName string) {
	delete(r.SkillList, skillName)
}
