package directions

import (
	. "interfaces"
)

type North struct{}

func (n North) Name() string {
	return north
}

func (n North) Execute(s SessionInterface, input []string) error {
	return execute(s, input, n.Name())
}
