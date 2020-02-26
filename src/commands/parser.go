package commands

import (
	"fmt"
)

type Parser struct {
	Queue chan string
}

func (p *Parser) Start() {
	p.Queue = make(chan string)
	defer close(p.Queue)

	for {
		select {
		case input := <-p.Queue:
			// It's probably important to NOT parse inputs concurrently to prevent possible race conditions
			err := p.parseInput(input)
			if err != nil {
				fmt.Printf("error parsing input: %v \n", err)
			}
			break
		}
	}
}

func (p *Parser) parseInput(input string) error {
	fmt.Printf("got input: \"%s\" \n", input)
	return nil
}
