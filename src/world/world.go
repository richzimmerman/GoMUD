package world

import (
	"container/list"
	"context"
	lib "lib/world"
	"logger"
)

var log = logger.NewLogger()

func LoadWorld(c context.Context) error {
	log.Info("Loading world.")
	// TODO: Read rooms from database, load into Room structs and map; Check for duplicate room IDs

	// Create a default limbo room so there's somewhere to go in an empty world
	limbo := &Room{
		id:          "0",
		zone:        "Limbo",
		description: "You're standing in limbo. There is nothing around you, and no where to go. You start to wonder, are you a god? Or is this some kind of terrible mistake?",
		players:     list.New(),
		nonPlayers:  list.New(),
	}
	// Ignoring the errror here... if limbo room already exists, cool
	lib.AddRoom(limbo)

	return nil
}
