package db

import (
	"github.com/andreafalzetti/comments-api/pkg/comments"
	"net"
)

type Connection interface {
	net.Conn
	GetClientId() string
}

type Record struct {
	ClientId        string
	IsAuthenticated bool
}

type State struct {
	clients     map[string]*Record
	discussions []*comments.Discussion
	connections map[string]Connection
}

// NewInMemoryDB creates a new in-memory key-value database used for development and testing
func NewInMemoryDB() *State {
	return &State{
		clients:     make(map[string]*Record),
		connections: make(map[string]Connection),
	}
}

// AddConnection adds a connection to the database
func (s *State) AddConnection(clientId string, conn Connection) {
	if s.connections[clientId] == nil {
		s.connections[clientId] = conn
	}
}

// AuthenticateClient authenticates a client
func (s *State) AuthenticateClient(clientId string) {
	r := s.clients[clientId]
	if r == nil {
		r = &Record{
			ClientId:        clientId,
			IsAuthenticated: false,
		}
		s.clients[clientId] = r
	}
	r.IsAuthenticated = true
}

func (s *State) GetClientById(clientId string) *Record {
	return s.clients[clientId]
}

func (s *State) GetDiscussionById(discussionId string) *comments.Discussion {
	for _, d := range s.discussions {
		if d.GetId() == discussionId {
			return d
		}
	}
	return nil
}

func (s *State) GetDiscussions() []*comments.Discussion {
	return s.discussions
}

func (s *State) AddDiscussion(discussion *comments.Discussion) {
	s.discussions = append(s.discussions, discussion)
}

// GetConnectionsByIds returns a list of connections by their ids
func (s *State) GetConnectionsByIds(clientIds []string) []Connection {
	connections := make([]Connection, 0)
	for _, clientId := range clientIds {
		connections = append(connections, s.connections[clientId])
	}
	return connections
}

// GetAuthenticatedConnections returns a list of authenticated connections
func (s *State) GetAuthenticatedConnections() []Connection {
	connections := make([]Connection, 0)
	for _, r := range s.clients {
		if r.IsAuthenticated {
			connections = append(connections, s.connections[r.ClientId])
		}
	}
	return connections
}

// GetClientIds returns a list of client ids
func (s *State) GetClientIds() []string {
	clientIds := make([]string, 0)
	for clientId := range s.clients {
		clientIds = append(clientIds, clientId)
	}
	return clientIds
}
