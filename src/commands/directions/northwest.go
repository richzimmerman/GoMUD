package directions

import (
	. "interfaces"
)

type Northwest struct{}

func (n Northwest) Execute(s SessionInterface, input []string) error {
	return execute(s, input, northwest)
}
