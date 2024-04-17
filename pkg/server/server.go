package server

import (
	"fmt"
	"github.com/andreafalzetti/comments-api/pkg/controller"
	"net"
)

// Server represents the server
type Server struct {
	host       string
	port       string
	controller *controller.Controller
}

// Config holds the configuration for the Server
type Config struct {
	Host       string
	Port       string
	Controller *controller.Controller
}

// New initializes a new instance of the Server
func New(config *Config) *Server {
	return &Server{
		host:       config.Host,
		port:       config.Port,
		controller: config.Controller,
	}
}

// Run starts the server
func (server *Server) Run() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", server.host, server.port))
	if err != nil {
		fmt.Printf("error listening for new connections: %v\n", err)
		return
	}
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			fmt.Printf("error closing connection: %v\n", err)
			return
		}
		fmt.Printf("connection closed successfully\n")
	}(listener)

	for {
		conn, err := listener.Accept()
		fmt.Printf("+++ new connection\n")

		if err != nil {
			fmt.Printf("error accepting new connection: %v\n", err)
		}

		server.controller.HandleNewConnection(conn)
	}
}
