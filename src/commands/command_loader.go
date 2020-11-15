package commands

import (
	. "interfaces"
	lib "lib/commands"
	"logger"
	"reflect"
	"strings"
)

var log = logger.NewLogger()

var commandRegistry = make(map[string]reflect.Type)

func LoadCommands() error {
	log.Info("Loading commands...")
	// Register a nil instance of each command type to expose them for use during runtime
	registerAllTypes()
	// Every type in the commandRegistry needs an instance created and stored in the Command map
	for key, _ := range commandRegistry {
		split := strings.Split(key, ".")
		commandName := split[len(split)-1]
		cmd := makeInstance(key)
		err := lib.AddCommand(commandName, cmd)
		if err != nil {
			log.Warn("failed to load command (%s): (%v)", commandName, err)
		}
	}
	return nil
}

func makeInstance(name string) CommandInterface {
	return reflect.New(commandRegistry[name]).Elem().Interface().(CommandInterface)
}

func registerType(typedNil interface{}) {
	t := reflect.TypeOf(typedNil).Elem()
	commandRegistry[t.PkgPath()+"."+t.Name()] = t
}

func registerAllTypes() {
	// This sucks and I don't like it, but Go doesn't do reflection quite like Java can so we can't dynamically load types
	// dynamically with reflection, we need to have a Type map. So (for now), we're registering a nil value of each command
	// type to look up when a user inputs a command
	for _, t := range nilCommandTypes {
		registerType(t)
	}
}
