package directions

import (
	. "interfaces"
)

type East struct{}

func (e East) Name() string {
	return east
}

func (e East) Execute(s SessionInterface, input []string) error {
	return execute(s, input, e.Name())
}
