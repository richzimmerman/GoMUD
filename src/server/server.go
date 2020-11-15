package server

import (
	"context"
	"logger"
	"net"
	"os"
	"os/signal"
	"syscall"
	"users/client"
	"utils"
)

const (
	port = ":4500"
	tcp  = "tcp4"
)

var log = logger.NewLogger()

// Storing this in a global variable to be able to reference outside of this package
// var Server *server

type server struct {
	Clients            map[string]*client.Client
	Listener           net.Listener
	disconnectListener chan string
}

func NewServer() (*server, error) {
	clients := make(map[string]*client.Client)
	listener, err := net.Listen(tcp, port)
	if err != nil {
		return nil, utils.Error(err)
	}
	return &server{
		Clients:            clients,
		Listener:           listener,
		disconnectListener: make(chan string),
	}, nil
}

func (s *server) handleConnection(c net.Conn) {
	newClient := client.NewClient(c)

	s.Clients[c.RemoteAddr().String()] = newClient
	log.Info("Connection received from: %s \n", c.RemoteAddr().String())

	err := newClient.GameLoop()
	if err != nil {
		log.Info("broken connection (%s): %v \n", c.RemoteAddr().String(), err)
	}
	s.disconnectListener <- c.RemoteAddr().String()
}

func makeOsSignalChannel() chan os.Signal {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	return sigs
}

func Start(ctx context.Context) {
	// Create a channel for interupt signals to cleanly shut down and save player data
	sigs := makeOsSignalChannel()

	context, cancel := context.WithCancel(ctx)
	defer cancel()

	Server, err := NewServer()
	if err != nil {
		log.Err("unable to start server: %v \n", err)
	}

	// Start a go routine to start accepting connections so this does not block the following for loop
	go func() {
		for {
			var c, err = Server.Listener.Accept()
			if err != nil {
				log.Err("unable to accept connection: %v \n", err)
			}
			go Server.handleConnection(c)
		}
	}()

	for {
		select {
		case <-sigs:
			cancel()
			break
		case <-context.Done():
			for _, client := range Server.Clients {
				log.Info("logging out connection at: %s", client.GetRemoteAddress())
				client.Logout()
			}
			return
		case disconnectAddress := <-Server.disconnectListener:
			log.Info("Disconnecting: %s \n", disconnectAddress)
			c, found := Server.Clients[disconnectAddress]
			if !found {
				log.Err("cannot find client for %s", disconnectAddress)
			}
			c.Logout()
			delete(Server.Clients, disconnectAddress)
			break
		}
	}
}
