package server

import (
	"context"
	"fmt"
	"net"
	"server/client"
	"utils"
)

const (
	PORT = ":4500"
)

// Storing this in a global variable to be able to reference outside of this package
var Server *server

type server struct {
	Clients            *map[string]*client.Client
	Listener           net.Listener
	disconnectListener chan string
}

func NewServer() (*server, error) {
	clients := make(map[string]*client.Client)
	listener, err := net.Listen("tcp4", PORT)
	if err != nil {
		return nil, utils.Error(err)
	}
	return &server{
		Clients: &clients,
		Listener: listener,
		disconnectListener: make(chan string),
	}, nil
}

func (s *server) handleConnection(c net.Conn) {
	newClient := client.NewClient(c)

	(*s.Clients)[c.RemoteAddr().String()] = newClient
	fmt.Printf("Connection received from: %s \n", c.RemoteAddr().String())

	err := newClient.GameLoop()
	if err != nil {
		fmt.Printf("broken connection (%s): %v \n", c.RemoteAddr().String(), err)
	}
	s.disconnectListener <- c.RemoteAddr().String()
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
