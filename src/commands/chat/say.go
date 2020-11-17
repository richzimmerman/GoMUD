package chat

import (
	"fmt"
	. "interfaces"
	lib "lib/world"
	"message"
	"strings"
	"utils"
)

type Say struct{}

func (sa Say) Execute(s SessionInterface, input []string) error {
	player := s.Player()
	chatMsg := strings.Join(input[1:], " ")
	firstPerson := fmt.Sprintf("You say, '%s'", chatMsg)
	thirdPerson := fmt.Sprintf("%s says, '%s'", player.GetName(), chatMsg)
	unformattedMsg := message.NewUnformattedMessage(firstPerson, "", thirdPerson)
	msg := message.NewMessage(player, nil, unformattedMsg)

	room, err := lib.GetRoom(player.GetLocation())
	if err != nil {
		return utils.Error(err)
	}
	room.Send(msg)
	return nil
}
