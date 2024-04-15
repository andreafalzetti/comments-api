package db

type connection struct {
	clientId      string
	authenticated bool
}

type State struct {
	connections map[string]connection
}

// NewInMemoryDB creates a new in-memory database used for development and testing
func NewInMemoryDB() *State {
	return &State{
		connections: make(map[string]connection),
	}
}

// AuthenticateClient authenticates a client
func (s *State) AuthenticateClient(clientId string) {
	s.connections[clientId] = connection{
		clientId:      clientId,
		authenticated: true,
	}
}

func (s *State) GetConnection(clientId string) connection {
	return s.connections[clientId]
}
