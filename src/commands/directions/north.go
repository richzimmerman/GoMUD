package directions

import (
	. "interfaces"
	lib "lib/world"
	"message"
	"utils"
)

type North struct{}

const north = "North"

var firstPerson = firstPersonMsg(north)
var thirdPeron = thirdPersonMsg(north)
var unformattedMsg = message.NewUnformattedMessage(firstPerson, "", thirdPeron)

func (n North) Name() string {
	return north
}

func (n North) Execute(s SessionInterface, input []string) error {
	player := s.Player()
	currentRoom, err := lib.GetRoom(player.GetLocation())
	if err != nil {
		return utils.Error(err)
	}
	exit, err := currentRoom.GetExit(north)
	if err != nil {
		s.Client().Out("You cannot go that way.")
		return nil
	}
	// TODO: check if sitting/laying and shit
	msg := message.NewMessage(player, nil, unformattedMsg)
	currentRoom.Send(msg)
	lib.MovePlayer(player, exit.Destination())
	return nil
}
