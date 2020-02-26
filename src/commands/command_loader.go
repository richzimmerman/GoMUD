package commands

import "reflect"

var CommandMap *map[string]reflect.Type


func LoadCommands() {
	c := make(map[string]reflect.Type)
	CommandMap = &c

	// TODO:
	// for command in reflect inspect package, get commands
	//    CommandMap[command.Name()] = command
}