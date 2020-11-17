package directions

import (
	. "interfaces"
)

type South struct{}

func (st South) Execute(s SessionInterface, input []string) error {
	return execute(s, input, south)
}
