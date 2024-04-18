package db

import (
	"github.com/andreafalzetti/comments-api/pkg/comments"
)

type Record struct {
	ClientId        string
	IsAuthenticated bool
}

type State struct {
	clients     map[string]*Record
	discussions []*comments.Discussion
}

// NewInMemoryDB creates a new in-memory key-value database used for development and testing
func NewInMemoryDB() *State {
	return &State{
		clients: make(map[string]*Record),
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
