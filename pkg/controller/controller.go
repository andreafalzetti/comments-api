package controller

import (
	"fmt"
	"github.com/andreafalzetti/comments-api/pkg/comments"
	"github.com/andreafalzetti/comments-api/pkg/db"
)

type Controller struct {
	db *db.State
}

func NewController(db *db.State) *Controller {
	return &Controller{
		db: db,
	}
}

func (c *Controller) HandleMessage(req *comments.Request) string {
	res := &comments.Response{}

	if req.ID != "" {
		res.ID = req.ID
	}

	// auth
	if req.Data == comments.ActionSignIn {
		c.db.AuthenticateClient(req.ClientID)
		return fmt.Sprintf("%s\n", req.ID)
	}

	// whoami
	if req.Data == comments.ActionWhoami {
		isAuth, clientId := c.db.IsAuthenticated(req.ClientID)
		if isAuth {
			return fmt.Sprintf("%s|%s\n", req.ID, clientId)
		} else {
			fmt.Println("NOT AUTHENTICATED")
		}
	}

	// signout
	if req.Data == comments.ActionSignOut {
		c.db.SignOut((req.ClientID))
		return fmt.Sprintf("%s\n", res.ID)
	}

	if req.ClientID != "" {
		currentState := c.db.GetConnection(req.ClientID)
		fmt.Println("state - ", currentState)
	}

	if res.Data != "" {
		return fmt.Sprintf("%s|%s\n", res.ID, res.Data)
	} else {
		return fmt.Sprintf("%s\n", res.ID)
	}
}
