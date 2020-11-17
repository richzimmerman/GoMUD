package chat

import (
	lib "lib/world"
	"testing"
	"tests"
	"world"

	"github.com/stretchr/testify/assert"
)

func TestSay(t *testing.T) {
	room := world.Limbo()
	lib.AddRoom(room)
	defer lib.RemoveRoom(room.Id())

	p1 := tests.NewMockPlayer()
	p1.SetName("Foo")

	p2 := tests.NewMockPlayer()
	p2.SetName("Bar")

	room.AddPlayer(p1)
	room.AddPlayer(p2)

	s := tests.NewMockSession(p1, nil)

	say := &Say{}
	err := say.Execute(s, []string{"say", "this", "is", "a", "test"})
	assert.Nil(t, err)

	p1Output := p1.GetOutput()
	assert.Equal(t, "You say, 'this is a test'", p1Output)
	p2Output := p2.GetOutput()
	assert.Equal(t, "Foo says, 'this is a test'", p2Output)
}

func TestSayErr(t *testing.T) {
	room := world.Limbo()
	p1 := tests.NewMockPlayer()

	room.AddPlayer(p1)
	s := tests.NewMockSession(p1, nil)

	say := &Say{}
	err := say.Execute(s, []string{"say", "this", "is", "a", "test"})
	assert.NotNil(t, err)
}
