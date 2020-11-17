package common

import (
	lib "lib/world"
	"testing"
	"tests"
	"world"

	"github.com/stretchr/testify/assert"
)

const expectedLook = `[<R>Limbo</R>]
You're standing in limbo. There is nothing around you, and no where to go. You start to wonder, are you a god? Or is this some kind of terrible mistake?

There are no obvious exits.
`

func TestLook(t *testing.T) {
	room := world.Limbo()
	lib.AddRoom(room)
	defer lib.RemoveRoom(room.Id())

	p := tests.NewMockPlayer()
	room.AddPlayer(p)

	s := tests.NewMockSession(p, nil)

	look := &Look{}
	err := look.Execute(s, nil)
	assert.Nil(t, err)

	output := p.GetOutput()
	assert.Equal(t, expectedLook, output)
}

func TestLookErr(t *testing.T) {
	room := world.Limbo()

	p := tests.NewMockPlayer()
	room.AddPlayer(p)

	s := tests.NewMockSession(p, nil)

	look := &Look{}
	err := look.Execute(s, nil)
	assert.NotNil(t, err)
}
