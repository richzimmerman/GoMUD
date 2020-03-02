package main

import (
	"context"
	"db"
	"fmt"
	"races"
	"server"
)

func main() {
	err := db.InitDatabaseConnection()
	if err != nil {
		panic(fmt.Errorf("unable to initialize db connection: %v", err))
	}

	defer db.DatabaseConnection.Connection.Close()

	// TODO: Load rooms/map from DB
	if err = races.LoadRaces(); err != nil {
		panic(fmt.Errorf("unable to load Races: %v", err))
	}

	// TODO: Load config for.. stuff?

	mainCtx := context.Background()

	server.Start(mainCtx)
}
