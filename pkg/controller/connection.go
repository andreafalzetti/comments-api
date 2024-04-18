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
		return fmt.Sprintf("%s|%s\n", req.ID, d.GetId())
	}

	// create reply
	if req.Data == comments.ActionCreateReply {
		fmt.Printf("New reply to discussion with ref: %s, saying: %s\n", req.Reference, req.Comment)
		d := conn.db.GetDiscussionById(req.DiscussionId)
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
		d := conn.db.GetDiscussionById(req.DiscussionId)
		if d == nil {
			fmt.Printf("x - error: discussion with ref %s not found\n", req.Reference)
			return ""
		}
		return fmt.Sprintf("%s|%s|%s|%s\n", req.ID, d.GetId(), d.GetReference(), d.GetReplies())
	}

	// list discussions
	if req.Data == comments.ActionListDiscussions {
		fmt.Printf("List discussions\n")
		ds := conn.db.GetDiscussions()
		var discussions []string
		for _, d := range ds {
			discussions = append(discussions, d.String())
		}
		return fmt.Sprintf("%s|(%s)\n", req.ID, strings.Join(discussions, ","))
	}

	// mcmbfkw|
	// (bb8ed992-afd8-401c-96ab-97bdc794cd4b|qsiqevw.0s|(Alpha|I love this video. What did you use to make it?,Bravo|"I used something called ""Synthesia"", it's pretty cool!")
	//  ,
	//  9d9f6504-e407-4ee7-b202-71b58b35cc12|qsiqevw.15s|(Alpha|I'm not sure about this title scene.,Charlie|"Yes, it's a bit too long, no?",Bravo|I'm not sure - I think it introduces the topic nicely.,Charlie|What about adding some animation to the scene elements?,Alpha|Or some music? Something to make it a little more interesting?,Bravo|I think these are good ideas. I'll see what I can do.))

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
