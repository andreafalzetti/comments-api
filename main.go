package main

import (
	"fmt"

	"github.com/andreafalzetti/comments-api/pkg/config"
	"github.com/andreafalzetti/comments-api/pkg/server"
)

// Main function
func main() {
	cfg, err := config.New()
	fmt.Println("Starting server...")
	s := server.New(&server.Config{Host: apiConfig.Host, Port: apiConfig.Host})
	s.Run()
}
