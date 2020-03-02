package races

import (
	"db"
	"fmt"
	"skills"
	"utils"
)

var Races map[string]*Race

type Race struct {
	Name        string
	Realm       int8
	Type        int8 // TODO: enum style int for race types (thief, crafty, spell caster, warrior, hybrid, tank) to effect skills/stats
	SkillList   map[string]*RaceSkill
	Description string
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
		return err
	}
	for _, race := range races {
		r := loadRace(race)
		Races[r.Name] = r
	}
	return nil
}

func loadRace(r *db.DBRace) *Race {
	// TODO: Skills
	// Unmarshal JSON string, iterate over skill list, load skill pointers into map
	// skills := ^that
	return &Race{
		Name:        r.Name,
		Realm:       r.Realm,
		Type:        r.Type,
		SkillList:   nil,
		Description: r.Description,
	}
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
	r := loadRace(d)
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
