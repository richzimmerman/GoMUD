package directions

import (
	. "interfaces"
)

type Southwest struct{}

func (st Southwest) Name() string {
	return southwest
}

func (st Southwest) Execute(s SessionInterface, input []string) error {
	return execute(s, input, st.Name())
}
