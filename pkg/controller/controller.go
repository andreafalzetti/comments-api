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
	if req.ClientID != "" {
		currentState := c.db.GetConnection(req.ClientID)
		fmt.Println("state - ", currentState)
	}

	if res.Data != "" {
		//return strings.Join(res.ID, res.Data, "|")
		return fmt.Sprintf("%s|%s", res.ID, res.Data)
	} else {
		return res.ID
	}
}
