package directions

import (
	. "interfaces"
)

type Out struct{}

func (o Out) Execute(s SessionInterface, input []string) error {
	return execute(s, input, out)
}
