package lib

import (
	"fmt"
	. "interfaces"
)

// Global store of rooms in the world
var rooms = make(map[string]RoomInterface)

func MoveMob(m MobInterface, roomId string) {
	// TODO: since we're handling NPCs differently, this should never be the case once spawners are implemented
	if m.Location() != "" {
		err := rooms[m.Location()].RemoveMob(m)
		if err != nil {
			fmt.Printf("%v", err)
		}
	}
	m.SetLocation(roomId)
	rooms[roomId].AddMob(m)
}

func MovePlayer(p PlayerInterface, roomId string) {
	if p.GetLocation() != "" {
		err := rooms[p.GetLocation()].RemovePlayer(p.GetDisplayName())
		if err != nil {
			// TODO: is this a problem? shouldn't be?
			fmt.Printf("%v\n", err)
		}
	}
	p.SetLocation(roomId)
	rooms[roomId].AddPlayer(p)
}

func AddRoom(r RoomInterface) error {
	if _, found := rooms[r.Id()]; found {
		return fmt.Errorf("room (%s) already exists", r.Id())
	}
	rooms[r.Id()] = r
	return nil
}

func RemoveRoom(roomId string) error {
	if _, found := rooms[roomId]; !found {
		return fmt.Errorf("room (%s) not found", roomId)
	}
	delete(rooms, roomId)
	return nil
}

func GetRoom(roomId string) (RoomInterface, error) {
	if _, found := rooms[roomId]; !found {
		return nil, fmt.Errorf("room (%s) not found", roomId)
	}
	return rooms[roomId], nil
}
