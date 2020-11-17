package common

import (
	. "interfaces"
	lib "lib/world"
	"utils"
)

type Look struct{}

func (c Look) Execute(s SessionInterface, input []string) error {
	p := s.Player()
	room, err := lib.GetRoom(p.GetLocation())
	if err != nil {
		return utils.Error(err)
	}
	p.Send(room.Look(p))
	return nil
}
