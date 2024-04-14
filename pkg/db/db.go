package db

type connection struct {
	clientId      int
	authenticated bool
}

type State struct {
	connections map[int]connection
}

// NewInMemoryDB creates a new in-memory database used for development and testing
func NewInMemoryDB() *State {
	return &State{
		connections: make(map[int]connection),
	}
}

// AuthenticateClient authenticates a client
func (s *State) AuthenticateClient(clientId int) {
	s.connections[clientId] = connection{
		clientId:      clientId,
		authenticated: true,
	}
}
