package directions

import (
	. "interfaces"
)

type Southeast struct{}

func (st Southeast) Execute(s SessionInterface, input []string) error {
	return execute(s, input, southeast)
}
