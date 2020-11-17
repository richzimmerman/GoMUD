package directions

import (
	lib "lib/world"
	"testing"
	"tests"
	"world"

	"github.com/stretchr/testify/assert"
)

func TestDirection(t *testing.T) {
	room0 := world.Limbo()
	lib.AddRoom(room0)
	defer lib.RemoveRoom(room0.Id())

	// For maximum code coverage!
	mockEast := tests.NewMockDirection("East", "1")
	mockUp := tests.NewMockDirection("Up", "1")
	mockDown := tests.NewMockDirection("Down", "1")
	mockOut := tests.NewMockDirection("Out", "1")
	room0.AddExit(mockEast)
	room0.AddExit(mockUp)
	room0.AddExit(mockDown)
	room0.AddExit(mockOut)

	room1 := world.Limbo()
	room1.SetId("1")
	lib.AddRoom(room1)
	defer lib.RemoveRoom(room1.Id())

	p1 := tests.NewMockPlayer()
	p1.SetName("Foo")
	p1.SetDisplayName("Foo")
	p2 := tests.NewMockPlayer()
	p2.SetName("Bar")
	p2.SetDisplayName("Bar")

	room0.AddPlayer(p1)
	room0.AddPlayer(p2)

	// TODO: Add player to room1 to get "X enterred the room" output
	s := tests.NewMockSession(p1, nil)

	err := execute(s, nil, "East")
	assert.Nil(t, err)

	p1Output := p1.GetOutput()
	assert.Equal(t, "You go east", p1Output)

	p2Output := p2.GetOutput()
	assert.Equal(t, "Foo leaves to the east", p2Output)

	lib.MovePlayer(p1, "0")

	err = execute(s, nil, "Up")
	p1Output = p1.GetOutput()
	assert.Equal(t, "You go up", p1Output)

	p2Output = p2.GetOutput()
	assert.Equal(t, "Foo goes up", p2Output)

	lib.MovePlayer(p1, "0")

	err = execute(s, nil, "Down")
	p1Output = p1.GetOutput()
	assert.Equal(t, "You go down", p1Output)

	p2Output = p2.GetOutput()
	assert.Equal(t, "Foo goes down", p2Output)

	lib.MovePlayer(p1, "0")

	err = execute(s, nil, "Out")
	p1Output = p1.GetOutput()
	assert.Equal(t, "You go out", p1Output)

	p2Output = p2.GetOutput()
	assert.Equal(t, "Foo goes out", p2Output)
}

func TestDirectionRoomNotFoundErr(t *testing.T) {
	p1 := tests.NewMockPlayer()

	s := tests.NewMockSession(p1, nil)

	err := execute(s, nil, "East")
	assert.NotNil(t, err)
}

func TestDirectionNotFoundErr(t *testing.T) {
	room := world.Limbo()
	lib.AddRoom(room)
	defer lib.RemoveRoom(room.Id())

	p1 := tests.NewMockPlayer()
	room.AddPlayer(p1)
	c := tests.NewMockClient()
	s := tests.NewMockSession(p1, c)

	err := execute(s, nil, "East")
	assert.Nil(t, err)

	output := c.GetOutput()
	assert.Equal(t, "You cannot go that way.", output)
}
