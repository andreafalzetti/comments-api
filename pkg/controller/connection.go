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
		fmt.Printf("<-- incoming request: '%s', '%s', '%s'\n", m.ID, m.Data, m.ClientID)
		output := conn.HandleMessage(m)

		conn.Write([]byte(output))
		//client.conn.Close()

		fmt.Printf("--> response: '%s'\n", strings.TrimSuffix(output, "\n"))
	}

}

func (conn *Connection) HandleMessage(req *comments.Request) string {
	res := &comments.Response{}

	if req.ID != "" {
		res.ID = req.ID
	}

	// auth
	if req.Data == comments.ActionSignIn {
		conn.db.AuthenticateClient(req.ClientID)
		conn.clientId = req.ClientID
		return fmt.Sprintf("%s\n", req.ID)
	}

	// whoami
	if req.Data == comments.ActionWhoami {
		fmt.Printf("DEBUG BEFORE IS AUTHENTICATED, CLIENT ID: %s\n", req)
		return fmt.Sprintf("%s|%s\n", req.ID, conn.WhoAmI())
	}

	// signout
	if req.Data == comments.ActionSignOut {
		conn.SignOut()
		return fmt.Sprintf("%s\n", res.ID)
	}

	// create discussion
	if req.Data == comments.ActionCreateDiscussion {
		fmt.Printf("New discussion with ref: %s, saying: %s\n", req.Reference, req.Comment)
		d := comments.NewDiscussion(req.Reference, conn.clientId, req.Comment)
		conn.db.AddDiscussion(d)
		return fmt.Sprintf("%s\n", d.GetId())
	}

	// create reply
	if req.Data == comments.ActionCreateReply {
		fmt.Printf("New reply to discussion with ref: %s, saying: %s\n", req.Reference, req.Comment)
		d := conn.db.GetDiscussionById(req.Reference)
		if d == nil {
			fmt.Printf("x - error: discussion with ref %s not found\n", req.Reference)
			return ""
		}
		d.AddReply(conn.clientId, req.Comment)
		return fmt.Sprintf("%s\n", req.ID)
	}

	// get discussion
	if req.Data == comments.ActionGetDiscussion {
		fmt.Printf("Get discussion with ref: %s\n", req.Reference)
		d := conn.db.GetDiscussionById(req.Reference)
		if d == nil {
			fmt.Printf("x - error: discussion with ref %s not found\n", req.Reference)
			return ""
		}
		return fmt.Sprintf("%s|%s\n", req.ID, d.GetId())
	}

	// list discussions
	if req.Data == comments.ActionListDiscussions {

	}

	if req.ClientID != "" {
		currentState := conn.db.GetClientById(req.ClientID)
		fmt.Println("db - ", currentState)
	}

	if res.Data != "" {
		return fmt.Sprintf("%s|%s\n", res.ID, res.Data)
	} else {
		return fmt.Sprintf("%s\n", res.ID)
	}
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
