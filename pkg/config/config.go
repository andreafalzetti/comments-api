package config

import "os"

type Config struct {
	Host string
	Port string
}

// New creates a new service config
func New(config *Config) (*Config, error) {

	host := os.Getenv("COMMENTS_API_HOST")

	if (host == "") {
		
	}

	return &Config{
		Host: ,
		Port: os.Getenv("COMMENTS_API_PORT"),
	}
}
