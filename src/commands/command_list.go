package commands

import (
	"commands/common"
	"commands/directions"
)

// For commands to be initialized at start up, a nil instance of each type needs to be added here
var nilCommandTypes = []interface{}{
	(*common.Look)(nil),
	(*directions.North)(nil),
	(*directions.Northeast)(nil),
	(*directions.Northwest)(nil),
	(*directions.South)(nil),
	(*directions.Southeast)(nil),
	(*directions.Southwest)(nil),
	(*directions.East)(nil),
	(*directions.West)(nil),
	(*directions.Out)(nil),
	(*directions.Up)(nil),
	(*directions.Down)(nil),
}
