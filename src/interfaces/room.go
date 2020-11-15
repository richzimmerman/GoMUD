package interfaces

type RoomInterface interface {
	Lock()
	Unlock()
	Id() string
	SetId(s string)
	Zone() string
	SetZone(zone string)
	Description() string
	SetDescription(description string)
	Exits() map[string]DirectionInterface
	GetExit(name string) (DirectionInterface, error)
	AddExit(dir DirectionInterface) error
	AddPlayer(player PlayerInterface)
	RemovePlayer(name string) error
	GetPlayer(name string) (PlayerInterface, error)
	AddMob(mob MobInterface)
	RemoveMob(mob MobInterface) error
	Look(self PlayerInterface) string
	Send(msg MessageInterface)
}
