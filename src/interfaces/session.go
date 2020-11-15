package interfaces

type SessionInterface interface {
	Player() PlayerInterface
	Client() ClientInterface
	LastInput() int64
	InputReceived(timestamp int64)
}
