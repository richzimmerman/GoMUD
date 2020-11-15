package interfaces

type MessageInterface interface {
	Antagonist() PlayerInterface
	Target() PlayerInterface
	UnformattedMessage() UnformatedMessageInterface
}

type UnformatedMessageInterface interface {
	FirstPerson() string
	SecondPerson() string
	ThirdPerson() string
}
