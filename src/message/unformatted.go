package message

type UnformattedMessage struct {
	firstPerson  string
	secondPerson string
	thirdPerson  string
}

func NewUnformattedMessage(f string, s string, t string) *UnformattedMessage {
	return &UnformattedMessage{
		firstPerson:  f,
		secondPerson: s,
		thirdPerson:  t,
	}
}

func (u *UnformattedMessage) FirstPerson() string {
	return u.firstPerson
}

func (u *UnformattedMessage) SecondPerson() string {
	return u.secondPerson
}

func (u *UnformattedMessage) ThirdPerson() string {
	return u.thirdPerson
}
