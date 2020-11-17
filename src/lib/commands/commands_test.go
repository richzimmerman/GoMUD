package commands

import (
	"fmt"
	. "interfaces"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockCommand struct{}

func (m *MockCommand) Execute(s SessionInterface, input []string) error {
	return fmt.Errorf("I ran a mock command")
}

func TestTrie(t *testing.T) {
	trie := initTrie()
	// Empty trie
	foo := trie.verify("foo")
	assert.False(t, foo)

	trie.insert("foo")

	fooTwo := trie.verify("foo")
	assert.True(t, fooTwo)
	barTwo := trie.verify("bar")
	assert.False(t, barTwo)

	trie.remove("foo")

	fooThree := trie.verify("foo")
	assert.False(t, fooThree)

	trie.insert("foo")

	fooFour := trie.verify("foo")
	assert.True(t, fooFour)
}

func TestCommandMap(t *testing.T) {
	fooCommand := &MockCommand{}

	foo := commandTrie.verify("foo")
	assert.False(t, foo)

	err := AddCommand("foo", fooCommand)
	assert.Nil(t, err)

	// trie does not normalize case of text since it's not a 'public' function
	fooTwo := commandTrie.verify("foo")
	assert.True(t, fooTwo)

	// Get command by full name
	cmd, err := GetCommand("foo")
	assert.Nil(t, err)

	assert.Equal(t, reflect.TypeOf(fooCommand), reflect.TypeOf(cmd))
	e := cmd.Execute(nil, nil)
	assert.Equal(t, "I ran a mock command", e.Error())

	err = RemoveCommand("foo")
	assert.Nil(t, err)
	check := commandTrie.verify("foo")
	assert.False(t, check)
	cmd, err = GetCommand("foo")
	assert.Nil(t, cmd)
	assert.Error(t, err)

	err = AddCommand("foobarbaz", fooCommand)
	assert.Nil(t, err)
	check = commandTrie.verify("foobarbaz")
	assert.True(t, check)

	// Get command by partial name
	cmdTwo, errTwo := GetCommand("fo")
	assert.Nil(t, errTwo)

	assert.Equal(t, reflect.TypeOf(fooCommand), reflect.TypeOf(cmdTwo))
	eTwo := cmdTwo.Execute(nil, nil)
	assert.Equal(t, "I ran a mock command", eTwo.Error())
}
