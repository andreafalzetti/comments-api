package db

import "fmt"

type Record struct {
	ClientId        string
	IsAuthenticated bool
}

type State struct {
	keys map[string]*Record
}

// NewInMemoryDB creates a new in-memory key-value database used for development and testing
func NewInMemoryDB() *State {
	return &State{
		keys: make(map[string]*Record),
	}
}

// AuthenticateClient authenticates a client
func (s *State) AuthenticateClient(clientId string) {
	fmt.Println("AuthenticateClient: ", clientId)
	r := s.keys[clientId]
	if r == nil {
		r = &Record{
			ClientId:        clientId,
			IsAuthenticated: false,
		}
		s.keys[clientId] = r
	}
	fmt.Println("Record: ", r)
	r.IsAuthenticated = true
}

func (s *State) GetRecordById(clientId string) *Record {
	return s.keys[clientId]
}

func (s *State) IsAuthenticated(clientId string) (bool, string) {
	fmt.Println("Client IsAuthenticated: ", clientId)
	r := s.keys[clientId]
	return r.IsAuthenticated, r.ClientId
}

func (s *State) SignOut(clientId string) {
	fmt.Println("Client SignOut: ", clientId)
	r := s.keys[clientId]
	r.IsAuthenticated = false
}
