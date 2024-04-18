package controller

import (
	"bufio"
	"fmt"
	"github.com/andreafalzetti/comments-api/pkg/comments"
	"github.com/andreafalzetti/comments-api/pkg/db"
	"net"
	"strings"
)

type Connection struct {
	net.Conn
	db       *db.State
	clientId string
}

func NewConnection(conn net.Conn, db *db.State) *Connection {
	return &Connection{
		Conn: conn,
		db:   db,
	}
}

// GetClientId returns the clientId of the authenticated client
func (conn *Connection) GetClientId() string {
	return conn.clientId
}

func (conn *Connection) Listen() {
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			if err.Error() != "EOF" {
				fmt.Printf("error: %v", err)
			}
			return
		}
		m := comments.Unmarshall(message)
		incomingClientId := conn.clientId
		if incomingClientId == "" {
			incomingClientId = conn.GetClientId()
		}
		fmt.Printf("<-- incoming request from: '%s': '%s'\n", incomingClientId, m.Data)
		outMsg, asyncOps := conn.HandleMessage(m)
		conn.db.AddConnection(conn.clientId, conn)
		_, err = conn.Write([]byte(outMsg))
		if err != nil {
			fmt.Printf("x - failed to write response to client: %s: '%s', with error: %v", conn.GetClientId(), strings.TrimSuffix(outMsg, "\n"), err)
			continue
		}

		fmt.Printf("--> reply to client '%s': '%s'\n", conn.GetClientId(), strings.TrimSuffix(outMsg, "\n"))

		if asyncOps != nil {
			fmt.Printf("--> running async ops\n")
			asyncOps()
		}
	}

}

func (conn *Connection) HandleMessage(req *comments.Request) (string, func()) {
	res := &comments.Response{}

	if req.ID != "" {
		res.ID = req.ID
	}

	// analyse the message to extract mentions
	mentions := []string{}
	if req.Comment != "" {
		allowedMentions := conn.db.GetClientIds()
		// check if comment contains "@{clientId}"
		for _, clientId := range allowedMentions {
			if strings.Contains(req.Comment, fmt.Sprintf("@%s", clientId)) {
				mentions = append(mentions, clientId)
			}
		}
	}

	// auth
	if req.Data == comments.ActionSignIn {
		conn.db.AuthenticateClient(req.ClientID)
		conn.clientId = req.ClientID
		return fmt.Sprintf("%s\n", req.ID), nil
	}

	// whoami
	if req.Data == comments.ActionWhoami {
		return fmt.Sprintf("%s|%s\n", req.ID, conn.WhoAmI()), nil
	}

	// signout
	if req.Data == comments.ActionSignOut {
		conn.SignOut()
		return fmt.Sprintf("%s\n", res.ID), nil
	}

	// create discussion
	if req.Data == comments.ActionCreateDiscussion {
		//fmt.Printf("New discussion with ref: %s, saying: %s\n", req.Reference, req.Comment)
		d := comments.NewDiscussion(req.Reference, conn.clientId, req.Comment)
		conn.db.AddDiscussion(d)
		return fmt.Sprintf("%s|%s\n", req.ID, d.GetId()), nil
	}

	// create reply
	if req.Data == comments.ActionCreateReply {
		//fmt.Printf("New reply to discussion with ref: %s, saying: %s\n", req.Reference, req.Comment)
		d := conn.db.GetDiscussionById(req.DiscussionId)
		if d == nil {
			fmt.Printf("x - error: discussion with ref %s not found\n", req.Reference)
			return "", nil
		}
		d.AddReply(conn.clientId, req.Comment)

		return fmt.Sprintf("%s\n", req.ID), func() {
			clientsToNotify := map[string]bool{}
			for _, p := range d.GetParticipants() {
				clientsToNotify[p] = true
			}
			for _, m := range mentions {
				clientsToNotify[m] = true
			}
			var clients []string
			for c := range clientsToNotify {
				clients = append(clients, c)
			}

			conns := conn.db.GetConnectionsByIds(clients)

			for _, c := range conns {
				if c.GetClientId() == conn.GetClientId() {
					fmt.Printf("--> skipping notification to client '%s' of discussion '%s'\n", c.GetClientId(), d.GetId())
					if c.GetClientId() == "Charlie" {
						_, err := c.Write([]byte(fmt.Sprintf("TEST\n")))
						if err != nil {
							fmt.Printf("x - error notifying participant: %v", err)
						}
						return
					}
					continue
				} else {
					fmt.Printf("--> notify client '%s' of discussion '%s'\n", c.GetClientId(), d.GetId())
					notification := fmt.Sprintf("%s|%s\n", comments.ActionDiscussionUpdated, d.GetId())
					_, err := c.Write([]byte(notification))
					if err != nil {
						fmt.Printf("x - error notifying participant: %v", err)
					}
				}
			}
		}
	}

	// get discussion
	if req.Data == comments.ActionGetDiscussion {
		fmt.Printf("Get discussion with ref: %s\n", req.Reference)
		d := conn.db.GetDiscussionById(req.DiscussionId)
		if d == nil {
			fmt.Printf("x - error: discussion with ref %s not found\n", req.Reference)
			return "", nil
		}
		return fmt.Sprintf("%s|%s|%s|%s\n", req.ID, d.GetId(), d.GetReference(), d.GetReplies()), nil
	}

	// list discussions
	if req.Data == comments.ActionListDiscussions {
		fmt.Printf("List discussions\n")
		ds := conn.db.GetDiscussions()
		var discussions []string
		for _, d := range ds {
			discussions = append(discussions, d.String())
		}
		return fmt.Sprintf("%s|(%s)\n", req.ID, strings.Join(discussions, ",")), nil
	}

	return "", nil
}

// WhoAmI handles the WHOAMI requests by returning the clientId of the authenticated client
func (conn *Connection) WhoAmI() string {
	r := conn.db.GetClientById(conn.clientId)
	if r != nil {
		if r.IsAuthenticated {
			return r.ClientId
		}
	}
	return ""
}

// SignOut handles the SIGN_OUT requests
func (conn *Connection) SignOut() {
	r := conn.db.GetClientById(conn.clientId)
	if r != nil {
		r.IsAuthenticated = false
	}
}
