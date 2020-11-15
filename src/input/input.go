package input

import "lib/sessions"

type Input struct {
	session *sessions.Session
	input   string
}

func (i *Input) Session() *sessions.Session {
	return i.session
}

func (i *Input) Input() string {
	return i.input
}

func NewInput(s *sessions.Session, i string) *Input {
	return &Input{
		session: s,
		input:   i,
	}
}
