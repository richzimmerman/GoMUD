package realms

import "fmt"

type Realm int

const (
	Evil  Realm = iota
	Chaos Realm = iota
	Good  Realm = iota
)

func (r Realm) String() string {
	switch r {
	case Evil:
		return "Evil"
	case Chaos:
		return "Chaos"
	case Good:
		return "Good"
	default:
		return fmt.Sprintf("%d", int(r))
	}
}

var Realms = []string{"Evil", "Chaos", "Good"}

// TODO: Read realms table from DB. create `Realm` with Name and ID number. change Realms to map[string]int or something
