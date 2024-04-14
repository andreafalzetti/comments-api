package controller

import (
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

func (c *Controller) HandleMessage(r *comments.Request) string {
	return r.ID
}
