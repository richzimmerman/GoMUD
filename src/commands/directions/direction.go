package directions

import (
	"fmt"
	. "interfaces"
	lib "lib/world"
	"message"
	"strings"
	"utils"
)

const (
	north                 = "North"
	northeast             = "Northeast"
	northwest             = "Northwest"
	south                 = "South"
	southeast             = "Southeast"
	southwest             = "Southwest"
	east                  = "East"
	west                  = "West"
	out                   = "Out"
	up                    = "Up"
	down                  = "Down"
	directionFirstPerson  = "You go %s"
	directionThirdPerson  = "<A.NAME> leaves to the %s"
	nonCardDirThirdPerson = "<A.NAME> goes %s"
)

func firstPersonMsg(dirName string) string {
	return fmt.Sprintf(directionFirstPerson, strings.ToLower(dirName))
}

func thirdPersonMsg(dirName string) string {
	var msg string
	switch dirName {
	case up, down, out:
		msg = fmt.Sprintf(nonCardDirThirdPerson, strings.ToLower(dirName))
		break
	default:
		msg = fmt.Sprintf(directionThirdPerson, strings.ToLower(dirName))
	}
	return msg
}

func execute(s SessionInterface, input []string, directionName string) error {
	firstPerson := firstPersonMsg(directionName)
	thirdPeron := thirdPersonMsg(directionName)
	unformattedMsg := message.NewUnformattedMessage(firstPerson, "", thirdPeron)

	player := s.Player()
	currentRoom, err := lib.GetRoom(player.GetLocation())
	if err != nil {
		return utils.Error(err)
	}
	exit, err := currentRoom.GetExit(directionName)
	if err != nil {
		s.Client().Out("You cannot go that way.")
		return nil
	}
	// TODO: check if sitting/laying and shit
	msg := message.NewMessage(player, nil, unformattedMsg)
	currentRoom.Send(msg)
	lib.MovePlayer(player, exit.Destination())
	// TODO: add "enterred" message in destination room, exit.Destination() maybe should be a RoomInterface
	return nil
}
