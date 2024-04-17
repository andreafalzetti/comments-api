package controller

import (
	"github.com/andreafalzetti/comments-api/pkg/db"
	"net"
)

type Controller struct {
	db *db.State
}

func NewController(db *db.State) *Controller {
	return &Controller{
		db: db,
	}
}

func (c *Controller) HandleNewConnection(netConn net.Conn) {
	conn := NewConnection(netConn, c.db)
	go conn.Listen()
}
