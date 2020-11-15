package world

type Direction struct {
	destination string
	name        string
}

func (d *Direction) Name() string {
	return d.name
}

func (d *Direction) SetName(name string) {
	d.name = name
}

func (d *Direction) Destination() string {
	return d.destination
}

func (d *Direction) SetDestination(roomId string) {
	d.destination = roomId
}
