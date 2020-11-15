package directions

import (
	. "interfaces"
)

type South struct{}

func (st South) Name() string {
	return south
}

func (st South) Execute(s SessionInterface, input []string) error {
	return execute(s, input, st.Name())
}
