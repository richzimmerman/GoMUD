package directions

import (
	. "interfaces"
)

type Northeast struct{}

func (n Northeast) Name() string {
	return northeast
}

func (n Northeast) Execute(s SessionInterface, input []string) error {
	return execute(s, input, n.Name())
}
