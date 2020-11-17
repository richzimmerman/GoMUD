package directions

import (
	. "interfaces"
)

type East struct{}

func (e East) Execute(s SessionInterface, input []string) error {
	return execute(s, input, east)
}
