package config

import (
	"fmt"
	"os"
)

type Config struct {
	Host string
	Port string
}

func loadEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// New creates a new service config
func NewFromEnv() (*Config, error) {
	host := loadEnvOrDefault("COMMENTS_API_HOST", "127.0.0.1")
	if host == "" {
		return nil, fmt.Errorf("host is missing")
	}

	port := loadEnvOrDefault("COMMENTS_API_PORT", "14000")
	if host == "" {
		return nil, fmt.Errorf("port is missing")
	}

	return &Config{
		Host: host,
		Port: port,
	}, nil
}
