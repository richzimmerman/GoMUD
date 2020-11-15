package directions

import (
	. "interfaces"
)

type Out struct{}

func (o Out) Name() string {
	return out
}

func (o Out) Execute(s SessionInterface, input []string) error {
	return execute(s, input, o.Name())
}
