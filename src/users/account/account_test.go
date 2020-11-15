package account

import (
	"db"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeleteCharacter(t *testing.T) {
	d := &db.DBAccount{}
	a := Account{
		DBAccount:  d,
		Characters: []string{"Foo", "Bar", "Baz"},
	}

	assert.Equal(t, a.Characters, []string{"Foo", "Bar", "Baz"})
	a.DeleteCharacter("Bar")
	assert.Equal(t, a.Characters, []string{"Foo", "Baz"})
	a.DeleteCharacter("Baz")
	assert.Equal(t, a.Characters, []string{"Foo"})
}
