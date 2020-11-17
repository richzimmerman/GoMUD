package directions

import (
	. "interfaces"
)

type Southwest struct{}

func (st Southwest) Execute(s SessionInterface, input []string) error {
	return execute(s, input, southwest)
}
