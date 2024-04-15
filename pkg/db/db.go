package db

import "fmt"

type Connection struct {
	clientId      string
	authenticated bool
}

type State struct {
	connections         map[string]Connection
	authenticatedClient string
}

// NewInMemoryDB creates a new in-memory database used for development and testing
func NewInMemoryDB() *State {
	return &State{
		connections: make(map[string]Connection),
	}
}

// AuthenticateClient authenticates a client
func (s *State) AuthenticateClient(clientId string) {
	fmt.Println("Client Authenticated: ", clientId)
	s.authenticatedClient = clientId
	s.connections[clientId] = Connection{
		clientId:      clientId,
		authenticated: true,
	}
}

func (s *State) GetConnection(clientId string) Connection {
	return s.connections[clientId]
}

func (s *State) IsAuthenticated(clientId string) (bool, string) {
	fmt.Println("Checking if client is authenticated: ", s.authenticatedClient)
	return s.authenticatedClient != "", s.authenticatedClient
	//return s.connections[clientId].authenticated, s.connections[clientId].clientId
}

func (s *State) SignOut(clientId string) {
	fmt.Println("Client Signed Out: ", clientId)
	s.authenticatedClient = ""
}
