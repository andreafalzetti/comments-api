package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

type Server struct {
	host string
	port string
}

type Client struct {
	conn net.Conn
}

type Config struct {
	Host string
	Port string
}

// New initializes a new instance of the Server
func New(config *Config) *Server {
	return &Server{
		host: config.Host,
		port: config.Port,
	}
}

// Run starts the server
func (server *Server) Run() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", server.host, server.port))
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		client := &Client{
			conn: conn,
		}
		go client.handleRequest()
	}
}

type message struct {
	requestId string // 7 chars (a-z), set by the client
	data      string
	clientId  string
}

func (message string) unmarshall() {
	// TODO: implement
}

func (client *Client) handleRequest() {
	reader := bufio.NewReader(client.conn)
	defer client.conn.Close()
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			if err.Error() != "EOF" {
				fmt.Printf("error: %v", err)
			}
			return
		}
		fmt.Printf("Message incoming: %s", string(message))
		client.conn.Write([]byte(message))
	}
}
