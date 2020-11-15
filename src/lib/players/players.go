package players

import (
	"db"
	"encoding/json"
	"fmt"
	. "interfaces"
	"races"
	"users/player"
	"utils"
	"world/realms"
)

var players = make(map[string]PlayerInterface)

// TODO: Load all the players from the DB at start up..
// these don't need to indicate LOGGED IN players, just all players

func AddPlayer(p PlayerInterface) error {
	if _, found := players[p.GetName()]; found {
		return fmt.Errorf("player (%s) already exists", p.GetName())
	}
	players[p.GetName()] = p
	return nil
}

func RemovePlayer(name string) error {
	if _, found := players[name]; !found {
		return fmt.Errorf("player (%s) not found", name)
	}
	delete(players, name)
	return nil
}

func GetPlayer(name string) (PlayerInterface, error) {
	if p, found := players[name]; found {
		return p, nil
	}
	return nil, fmt.Errorf("player (%s) not found", name)
}

func CheckPlayer(name string, c ClientInterface) (PlayerInterface, error) {
	result, err := db.DatabaseConnection.QueryPlayer(name)
	if err != nil {
		return nil, err
	}
	player, err := LoadPlayer(c, result)
	if err != nil {
		return nil, err
	}
	return player, nil
}

// NewPlayer generates a new player struct and saves it to the database. Really only used during character creation
func NewPlayer(c ClientInterface, account string, name string, race string, realm int8) (PlayerInterface, error) {
	r, ok := races.Races[race]
	if !ok {
		return nil, utils.Error(fmt.Errorf("invalid race, or race not loaded: %s", race))
	}
	if r.Realm != realm {
		return nil, utils.Error(fmt.Errorf("realm mismatch, %s is not part of realm %s",
			race, realms.Realms[realm]))
	}

	p := &player.Player{
		Name:        name,
		DisplayName: name,
		Level:       1,
		Health:      r.DefaultHealth,
		Fatigue:     r.DefaultFatigue,
		Power:       r.DefaultPower,
		Buffs:       nil,
		Debuffs:     nil,
		Location:    r.StartingRoom,
		Client:      c,
		Account:     account,
		Title:       r.DefaultTitle,
		RealmTitle:  "",
		Race:        r,
		PlayerStats: r.DefaultStats,
		Stance:      0,
		Skills:      nil,
		Spells:      nil,
	}
	err := SyncPlayer(p)
	if err != nil {
		return nil, utils.Error(err)
	}
	return p, nil
}

// SyncPlayer sync's the players data to the database, saving it's state
func SyncPlayer(p *player.Player) error {
	/* Converts a *Player to a *DBPlayer and saves it in the DB. */
	var pStats string
	s, err := json.Marshal(p.PlayerStats)
	if err != nil {
		return utils.Error(err)
	}
	pStats = string(s)

	sk, err := utils.JsonMarshalMapToArray(p.Skills)
	if err != nil {
		return utils.Error(err)
	}

	sp, err := utils.JsonMarshalMapToArray(p.Spells)
	if err != nil {
		return utils.Error(err)
	}

	bf, err := json.Marshal(p.Buffs)
	if err != nil {
		return utils.Error(err)
	}
	buffs := string(bf)
	if buffs == "null" {
		buffs = "{}"
	}

	dbf, err := json.Marshal(p.Debuffs)
	if err != nil {
		return utils.Error(err)
	}
	debuffs := string(dbf)
	if debuffs == "null" {
		debuffs = "{}"
	}

	d := &db.DBPlayer{
		Name:        p.GetName(),
		Account:     p.AccountName(),
		DisplayName: p.GetDisplayName(),
		Level:       p.GetLevel(),
		Health:      p.GetHealth(),
		Fatigue:     p.GetFatigue(),
		Power:       p.GetPower(),
		Title:       p.GetTitle(),
		RealmTitle:  p.GetRealmTitle(),
		Race:        p.RaceName(),
		Stats:       pStats,
		Stance:      0,
		Skills:      sk,
		Spells:      sp,
		Buffs:       buffs,
		Debuffs:     debuffs,
		Location:    p.GetLocation(),
	}

	err = db.DatabaseConnection.SavePlayer(d)
	return utils.Error(err)
}

// LoadPlayer loads a player struct from a struct representation of the data stored in the database
func LoadPlayer(c ClientInterface, p *db.DBPlayer) (PlayerInterface, error) {
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

	// TODO: parse JSON for stats, skills, spells
	race, ok := races.Races[p.Race]
	if !ok {
		return nil, utils.Error(fmt.Errorf("invalid race, or race not loaded: %s", p.Race))
	}
	return &player.Player{
		Name:        p.Name,
		DisplayName: p.DisplayName,
		Level:       p.Level,
		Health:      p.Health,
		Fatigue:     p.Fatigue,
		Power:       p.Power,
		Buffs:       buffs,
		Debuffs:     debuffs,
		Location:    p.Location,
		Client:      c,
		Account:     p.Account,
		Title:       p.Title,
		RealmTitle:  p.RealmTitle,
		Race:        race,
		PlayerStats: stats,
		Stance:      0,
		Skills:      nil,
		Spells:      nil,
	}, nil
}

func GetPlayersForAccount(accountName string) []string {
	p := make([]string, 0)
	for name, player := range players {
		if player.AccountName() == accountName {
			p = append(p, name)
		}
	}
	return p
}
