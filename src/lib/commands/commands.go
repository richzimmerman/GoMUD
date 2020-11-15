package commands

import (
	"fmt"
	. "interfaces"
	"logger"
	"strings"
)

var log = logger.NewLogger()

const (
	//ALBHABET_SIZE total characters in english alphabet
	ALBHABET_SIZE = 26
)

var commandTrie = initTrie()
var commandMap = make(map[string]CommandInterface)

func AddCommand(name string, c CommandInterface) error {
	cmdName := strings.ToLower(name)
	if _, found := commandMap[cmdName]; found {
		return fmt.Errorf("command (%s) already exists", name)
	}
	commandMap[cmdName] = c
	commandTrie.insert(cmdName)
	return nil
}

func RemoveCommand(name string) error {
	cmdName := strings.ToLower(name)
	if err := commandTrie.remove(cmdName); err != nil {
		log.Err("command (%s) is not indexed", name)
	}
	if _, found := commandMap[cmdName]; !found {
		return fmt.Errorf("command (%s) not found", cmdName)
	}
	delete(commandMap, cmdName)
	return nil
}

func GetCommand(name string) (CommandInterface, error) {
	cmdName := strings.ToLower(name)
	return commandTrie.find(cmdName)
}

func findCommand(name string) (CommandInterface, error) {
	cmd, found := commandMap[name]
	if !found {
		return nil, fmt.Errorf("command (%s) not found", name)
	}
	return cmd, nil
}

type trieNode struct {
	childrens [ALBHABET_SIZE]*trieNode
	command   CommandInterface
	word      string
	isWordEnd bool
}

type trie struct {
	root *trieNode
}

func initTrie() *trie {
	return &trie{
		root: &trieNode{},
	}
}

func (t *trie) insert(word string) {
	word = strings.ToLower(word)
	wordLength := len(word)
	current := t.root
	for i := 0; i < wordLength; i++ {
		index := word[i] - 'a'
		if current.childrens[index] == nil {
			current.childrens[index] = &trieNode{}
		}
		current = current.childrens[index]
	}
	current.word = strings.ToLower(word)
	current.command = commandMap[current.word]
	current.isWordEnd = true
}

func (t *trie) verify(word string) bool {
	word = strings.ToLower(word)
	wordLength := len(word)
	current := t.root
	for i := 0; i < wordLength; i++ {
		index := word[i] - 'a'
		if current.childrens[index] == nil {
			return false
		}
		current = current.childrens[index]
	}
	if current.isWordEnd {
		return true
	}
	return false
}

func (t *trie) find(word string) (CommandInterface, error) {
	word = strings.ToLower(word)
	wordLength := len(word)
	current := t.root
	for i := 0; i < wordLength; i++ {
		index := word[i] - 'a'
		if current.childrens[index] == nil {
			return nil, fmt.Errorf("command not found")
		}
		current = current.childrens[index]
	}
	if current.isWordEnd {
		return current.command, nil // return command, nil as this should guarantee an existing command name
	} else {
		return t.getFirstChild(current), nil
	}
}

func (t *trie) remove(word string) error {
	word = strings.ToLower(word)
	if !t.verify(word) {
		return fmt.Errorf("command (%s) does not exist", word)
	}
	wordLength := len(word)
	current := t.root
	var currentParent *trieNode
	var index byte
	for i := 0; i < wordLength; i++ {
		index = word[i] - 'a'
		c := current.childrens[index]
		currentParent = current
		current = c
	}
	currentParent.childrens[index] = nil
	return nil
}

func (t *trie) getFirstChild(node *trieNode) CommandInterface {
	for _, value := range node.childrens {
		if value != nil {
			if value.isWordEnd {
				return value.command
			} else {
				return t.getFirstChild(value)
			}
		}
	}
	// This should never happen since `getFirstChild` is only called when we've found a starting node
	return nil
}
