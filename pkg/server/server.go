package server

import (
	"bufio"
	"fmt"
	"github.com/andreafalzetti/comments-api/pkg/comments"
	"github.com/andreafalzetti/comments-api/pkg/controller"
	"log"
	"net"
	"strings"
)

type Server struct {
	host           string
	port           string
	requestHandler *controller.Controller
}

type Client struct {
	conn           net.Conn
	requestHandler *controller.Controller
}

type Config struct {
	Host           string
	Port           string
	RequestHandler *controller.Controller
}

// New initializes a new instance of the Server
func New(config *Config) *Server {
	return &Server{
		host:           config.Host,
		port:           config.Port,
		requestHandler: config.RequestHandler,
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
			conn:           conn,
			requestHandler: server.requestHandler,
		}
		go client.handleRequest()
	}
}
func unmarshall(raw string) *comments.Request {
	// the message is divided into parts separated by pipe |
	// first 7 chars are the request id
	// the following segment is the data
	// the last is the client id (optional)
	r := &comments.Request{}
	raw = strings.TrimSuffix(raw, "\n")
	parts := strings.Split(raw, "|")

	r.ID = parts[0]
	if len(parts) == 3 {
		r.Data = parts[1]
		r.ClientID = parts[2]
	} else if len(parts) == 2 {
		r.Data = parts[1]
	}
	return r
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
		m := unmarshall(message)
		fmt.Printf("<-- incoming request: '%s' - '%s' - '%s'\n", m.ID, m.Data, m.ClientID)
		output := client.requestHandler.HandleMessage(m)

		write, err := client.conn.Write([]byte(output))
		if err != nil {
			fmt.Printf("error: %v", err)
		}
		fmt.Printf("--> response: '%s' - %d\n", output, write)
	}
}
