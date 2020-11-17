package directions

import (
	. "interfaces"
)

type North struct{}

func (n North) Execute(s SessionInterface, input []string) error {
	return execute(s, input, north)
}
