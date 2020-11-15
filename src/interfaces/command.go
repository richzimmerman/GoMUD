package interfaces

type CommandInterface interface {
	Execute(s SessionInterface, input []string) error
	Name() string
}
