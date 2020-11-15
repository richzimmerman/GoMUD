package directions

import (
	. "interfaces"
)

type Southeast struct{}

func (st Southeast) Name() string {
	return southeast
}

func (st Southeast) Execute(s SessionInterface, input []string) error {
	return execute(s, input, st.Name())
}
