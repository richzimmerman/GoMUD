package message

import (
	. "interfaces"
	"strings"
)

type Perspective int

const (
	FirstPerson  Perspective = iota
	SecondPerson Perspective = iota
	ThirdPerson  Perspective = iota

	antagonist = "<A.NAME>"
	target     = "<T.NAME>"
)

func MessageFormatter(msg MessageInterface, p Perspective) string {
	var s string
	switch p {
	case FirstPerson:
		s = msg.UnformattedMessage().FirstPerson()
		if msg.Target() != nil {
			s = strings.ReplaceAll(s, target, msg.Target().GetName())
		}
		break
	case SecondPerson:
		s = msg.UnformattedMessage().SecondPerson()
		if msg.Antagonist() != nil {
			s = strings.ReplaceAll(s, antagonist, msg.Antagonist().GetName())
		}
		break
	default:
		s = msg.UnformattedMessage().ThirdPerson()
		if msg.Antagonist() != nil {
			s = strings.ReplaceAll(s, antagonist, msg.Antagonist().GetName())
		}
		if msg.Target() != nil {
			s = strings.ReplaceAll(s, target, msg.Target().GetName())
		}
	}
	return s
}

type Message struct {
	antagonist         PlayerInterface
	target             PlayerInterface
	unformattedMessage UnformatedMessageInterface
}

func NewMessage(a PlayerInterface, t PlayerInterface, msg UnformatedMessageInterface) MessageInterface {
	return &Message{
		antagonist:         a,
		target:             t,
		unformattedMessage: msg,
	}
}

func (m *Message) Antagonist() PlayerInterface {
	return m.antagonist
}

func (m *Message) Target() PlayerInterface {
	return m.target
}

func (m *Message) UnformattedMessage() UnformatedMessageInterface {
	return m.unformattedMessage
}
