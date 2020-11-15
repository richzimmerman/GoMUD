package interfaces

type ClientInterface interface {
	GetRemoteAddress() string
	Logout()
	GameLoop() error
	Out(msg string)
}
