package server

import (
	"context"
	"fmt"
	"net"
)

const (
	PORT = ":4500"
)

// Storing this in a global variable to be able to reference outside of this package
var Server *server

type server struct {
	Clients            *map[string]*Client
	Listener           net.Listener
	disconnectListener chan string
}

func NewServer() (*server, error) {
	clients := make(map[string]*Client)
	listener, err := net.Listen("tcp4", PORT)
	if err != nil {
		return nil, err
	}
	return &server{
		Clients: &clients,
		Listener: listener,
		disconnectListener: make(chan string),
	}, nil
}

func (s *server) handleConnection(c net.Conn) {
	client := NewClient(c)

	(*s.Clients)[client.Address] = client
	fmt.Printf("Connection received from: %s \n", client.Address)

	err := client.gameLoop()
	if err != nil {
		fmt.Printf("broken connection (%s): %v \n", c.RemoteAddr().String(), err)
	}
	s.disconnectListener <- client.Address
}

func Start(ctx context.Context) {
	var err error
	Server, err = NewServer()
	if err != nil {
		fmt.Printf("unable to start server: %v \n", err)
	}

	// Start a go routine to start accepting connections so this does not block the following for loop
	go func() {
		for {
			var c, err = Server.Listener.Accept()
			if err != nil {
				fmt.Printf("unable to accept connection: %v \n", err)
			}
			go Server.handleConnection(c)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case disconnectAddress := <-Server.disconnectListener:
			fmt.Printf("Disconnecting: %s \n", disconnectAddress)
			delete(*Server.Clients, disconnectAddress)
			break
		}
	}
}
