package commands

import (
	"commands/common"
	"commands/directions"
)

// For commands to be initialized at start up, a nil instance of each type needs to be added here
var nilCommandTypes = []interface{}{
	(*common.Look)(nil),
	(*directions.North)(nil),
}
