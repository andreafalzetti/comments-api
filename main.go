package main

import (
	"fmt"
	"github.com/andreafalzetti/comments-api/pkg/config"
	"github.com/andreafalzetti/comments-api/pkg/controller"
	"github.com/andreafalzetti/comments-api/pkg/db"
	"github.com/andreafalzetti/comments-api/pkg/server"
)

func main() {
	appCfg, err := config.FromEnv()
	if err != nil {
		fmt.Printf("error: %s", err)
		return
	}

	// create an in-memory state for the application
	inMemoryDB := db.NewInMemoryDB()

	// create a new TCP server instance
	serverCfg := &server.Config{
		Host:       appCfg.Host,
		Port:       appCfg.Port,
		Controller: controller.NewController(inMemoryDB),
	}
	s := server.New(serverCfg)

	// start the service
	s.Run()

	fmt.Printf("Starting server on port %s...\n", appCfg.Port)
}
