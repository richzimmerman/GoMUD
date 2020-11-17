package directions

import (
	. "interfaces"
)

type Down struct{}

func (d Down) Execute(s SessionInterface, input []string) error {
	return execute(s, input, down)
}
