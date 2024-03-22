package main

import (
	"fmt"

	"github.com/andreafalzetti/comments-api/pkg/config"
	"github.com/andreafalzetti/comments-api/pkg/server"
)

// Main function
func main() {
	cfg, err := config.NewFromEnv()
	if err != nil {
		fmt.Printf("error: %w", err)
		return
	}
	fmt.Printf("Starting server on port %s...", cfg.Port)
	s := server.New(&server.Config{Host: cfg.Host, Port: cfg.Host})
	s.Run()
}
