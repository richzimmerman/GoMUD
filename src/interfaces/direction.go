package interfaces

type DirectionInterface interface {
	Name() string
	SetName(name string)
	Destination() string
	SetDestination(roomId string)
}
