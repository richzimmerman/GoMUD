package commands


type command interface {
	Execute(...interface{}) error
	Name() string
}
