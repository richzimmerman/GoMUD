package world

import (
	"container/list"
	"fmt"
	"interfaces"
	. "interfaces"
	"message"
	"sort"
	"strings"
	"sync"
)

type Room struct {
	lock        sync.Mutex
	id          string
	zone        string // TODO: Maybe make this an ID
	description string
	exits       map[string]DirectionInterface
	players     *list.List
	nonPlayers  *list.List
	// TODO: Items, Structures
	CommandQueue chan string
}

func (r *Room) Lock() {
	r.lock.Lock()
}

func (r *Room) Unlock() {
	r.lock.Unlock()
}

func (r *Room) Id() string {
	return r.id
}

func (r *Room) SetId(id string) {
	r.id = id
}

func (r *Room) Zone() string {
	return r.zone
}

func (r *Room) SetZone(zone string) {
	r.Lock()
	defer r.Unlock()

	r.zone = zone
}

func (r *Room) Description() string {
	return r.description
}

func (r *Room) SetDescription(description string) {
	r.Lock()
	defer r.Unlock()

	r.description = description
}

func (r *Room) Exits() map[string]DirectionInterface {
	return r.exits
}

func (r *Room) getExitList() []string {
	e := []string{}
	for key, _ := range r.exits {
		e = append(e, key)
	}
	sort.Strings(e)
	return e
}

func (r *Room) GetExit(name string) (DirectionInterface, error) {
	if e, found := r.exits[name]; found {
		return e, nil
	}
	return nil, fmt.Errorf("exit (%s) not found", name)
}

func (r *Room) AddExit(dir DirectionInterface) error {
	r.Lock()
	defer r.Unlock()

	if _, found := r.exits[dir.Name()]; found {
		return fmt.Errorf("exit (%s) already exists", dir.Name())
	}
	r.exits[dir.Name()] = dir
	return nil
}

func (r *Room) AddPlayer(player PlayerInterface) {
	r.Lock()
	defer r.Unlock()

	r.players.PushBack(player)
}

func (r *Room) RemovePlayer(name string) error {
	r.Lock()
	defer r.Unlock()
	// TODO: Probably want to check for equals value first before accepting prefix for better targetting
	player := r.getPlayerElement(name)
	if player == nil {
		return fmt.Errorf("player (%s) is not in room (%s)", name, r.id)
	}
	r.players.Remove(player)
	return nil
}

func (r *Room) GetPlayer(name string) (PlayerInterface, error) {
	p := r.getPlayer(name)
	if p == nil {
		return nil, fmt.Errorf("player (%s) not found in room", name)
	}
	return p, nil
}

func (r *Room) getPlayerElement(name string) *list.Element {
	/*
		This helper function returns the *list.Element for the first player matching provided input
		for modifying the List
	*/
	for e := r.players.Front(); e != nil; e = e.Next() {
		lower := strings.ToLower(e.Value.(PlayerInterface).GetDisplayName())
		if strings.HasPrefix(lower, strings.ToLower(name)) {
			return e
		}
	}
	return nil
}

func elementToPlayer(e *list.Element) PlayerInterface {
	return e.Value.(interfaces.PlayerInterface)
}

func (r *Room) getPlayer(name string) PlayerInterface {
	element := r.getPlayerElement(name)
	if element != nil {
		return elementToPlayer(element)
	}
	return nil
}

func (r *Room) getFriendlyPlayers(self PlayerInterface) []string {
	s := []string{}
	for e := r.players.Front(); e != nil; e = e.Next() {
		player := elementToPlayer(e)
		if player.GetDisplayName() == self.GetDisplayName() {
			continue
		} else {
			if player.Realm() == self.Realm() {
				s = append(s, elementToPlayer(e).GetDisplayName())
			}
		}
	}
	return s
}

func (r *Room) getEnemyPlayers(self PlayerInterface) []string {
	s := []string{}
	for e := r.players.Front(); e != nil; e = e.Next() {
		player := elementToPlayer(e)
		if player.Realm() != self.Realm() {
			s = append(s, elementToPlayer(e).GetDisplayName())
		}
	}
	return s
}

func (r *Room) AddMob(mob MobInterface) {
	r.Lock()
	defer r.Unlock()

	r.nonPlayers.PushBack(mob)
}

func (r *Room) RemoveMob(mob MobInterface) error {
	r.Lock()
	defer r.Unlock()

	e := r.getNPCElement(mob)
	if e == nil {
		return fmt.Errorf("mob (%s, GUID: %s) does not exist in room (%s)", mob.DisplayName(),
			mob.GetGUID(), r.id)
	}
	r.nonPlayers.Remove(e)
	return nil
}

func (r *Room) getNPCElement(mob MobInterface) *list.Element {
	for e := r.nonPlayers.Front(); e != nil; e = e.Next() {
		guid := strings.ToLower(e.Value.(MobInterface).GetGUID())
		if mob.GetGUID() == guid {
			return e
		}
	}
	return nil
}

func (r *Room) Look(self PlayerInterface) string {
	output := fmt.Sprintf("[<R>%s</R>]\n%s\n\n", r.Zone(), r.Description()) // Zone Name and room Description

	// Add exits
	exits := r.getExitList()
	if len(r.exits) > 0 {
		output += "Obvious exists: "
		for i, exit := range exits {
			if i == len(exits)-1 {
				output += fmt.Sprintf("and %s.", strings.ToLower(exit))
			} else {
				output += fmt.Sprintf("%s, ", strings.ToLower(exit))
			}
		}
	} else {
		output += "There are no obvious exits."
	}
	output += "\n"

	// Add players in room, if any
	players := r.getFriendlyPlayers(self)
	// TODO: combine players and friendly npcs to list and iterate over that instead of `players`
	if len(players) > 0 {
		if len(players) > 1 {
			for i, p := range players {
				if i == len(players)-1 {
					output += fmt.Sprintf("and <G>%s</G> are also here.\n", p)
				} else {
					output += fmt.Sprintf("<G>%s</G>, ", p)
				}
			}
		} else {
			output += fmt.Sprintf("<G>%s</G> is also here.\n", players[0])
		}
	}
	// TODO: Get enemy npcs, combine list with enemy players
	enemies := r.getEnemyPlayers(self)
	if len(enemies) > 0 {
		output += fmt.Sprintf("Also there ")
		if len(enemies) > 1 {
			for i, e := range enemies {
				if i == len(enemies)-1 {
					output += fmt.Sprintf("and <R>%s</R>.\n", e)
				} else {
					output += fmt.Sprintf("<R>%s</R>, ", e)
				}
			}
		} else {
			output += fmt.Sprintf("<R>%s</R>.\n", enemies[0])
		}
	}
	// TODO: Add items on ground
	return output
}

func (r *Room) Send(msg MessageInterface) {
	/*
		Generic message that can be formatted the same for all recipients
	*/
	r.Lock()
	defer r.Unlock()

	for e := r.players.Front(); e != nil; e = e.Next() {
		p, ok := e.Value.(PlayerInterface)
		if !ok {
			continue
		}
		perspective := message.ThirdPerson
		playerName := p.GetName()
		if msg.Antagonist() != nil && msg.Antagonist().GetName() == playerName {
			perspective = message.FirstPerson
		}
		if msg.Target() != nil && msg.Target().GetName() == playerName {
			perspective = message.SecondPerson
		}
		m := message.MessageFormatter(msg, perspective)
		e.Value.(PlayerInterface).Send(m)
	}
}
