package directions

import (
	. "interfaces"
)

type Up struct{}

func (u Up) Execute(s SessionInterface, input []string) error {
	return execute(s, input, up)
}
