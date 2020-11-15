package main

import (
	"commands"
	"config"
	"context"
	"db"
	"fmt"
	. "lib/accounts"
	"logger"
	"races"
	"server"
	"world"
)

const (
	configFile = "Config"
)

var log = logger.NewLogger()

func main() {
	// TODO: Load in a configuration file
	err := config.LoadConfiguration(configFile)
	if err != nil {
		panic(err)
	}
	// TODO: set up log file
	//    defer file.close()
	mainCtx := context.Background()

	err = db.InitDatabaseConnection()
	if err != nil {
		panic(fmt.Errorf("unable to initialize db connection: %v", err))
	}

	defer db.DatabaseConnection.Connection.Close()

	// TODO: Load config for.. stuff?
	if err = world.LoadWorld(mainCtx); err != nil {
		panic(fmt.Errorf("unable to load World: %v", err))
	}

	// TODO: Load rooms/map from DB
	if err = races.LoadRaces(); err != nil {
		panic(fmt.Errorf("unable to load Races: %v", err))
	}

	if err = commands.LoadCommands(); err != nil {
		panic(err)
	}

	if err = LoadAccounts(); err != nil {
		panic(err)
	}

	log.Info("Server started.")
	server.Start(mainCtx)
}
