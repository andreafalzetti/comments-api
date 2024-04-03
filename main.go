package main

import (
	"fmt"

	"github.com/andreafalzetti/comments-api/pkg/config"
	"github.com/andreafalzetti/comments-api/pkg/server"
)

// Main function
func main() {
	cfg, err := config.FromEnv()
	if err != nil {
		fmt.Printf("error: %s", err)
		return
	}
	fmt.Printf("Starting server on port %s...\n", cfg.Port)
	s := server.New(&server.Config{Host: cfg.Host, Port: cfg.Port})
	s.Run()
}
