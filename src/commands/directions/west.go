package directions

import (
	. "interfaces"
)

type West struct{}

func (w West) Name() string {
	return south
}

func (w West) Execute(s SessionInterface, input []string) error {
	return execute(s, input, w.Name())
}
