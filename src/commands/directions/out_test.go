package directions

import (
	lib "lib/world"
	"testing"
	"tests"
	"world"

	"github.com/stretchr/testify/assert"
)

func TestOut(t *testing.T) {
	room0 := world.Limbo()
	lib.AddRoom(room0)
	defer lib.RemoveRoom(room0.Id())

	mockExit := tests.NewMockDirection("Out", "1")
	room0.AddExit(mockExit)

	room1 := world.Limbo()
	room1.SetId("1")
	lib.AddRoom(room1)
	defer lib.RemoveRoom(room1.Id())

	p1 := tests.NewMockPlayer()
	p1.SetName("Foo")
	p1.SetDisplayName("Foo")

	room0.AddPlayer(p1)

	// TODO: Add player to room1 to get "X enterred the room" output
	s := tests.NewMockSession(p1, nil)
	d := &Out{}
	err := d.Execute(s, nil)
	assert.Nil(t, err)

	p1Output := p1.GetOutput()
	assert.Equal(t, "You go out", p1Output)
}
