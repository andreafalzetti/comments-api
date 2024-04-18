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

// FromEnv creates a new service config reading values from the env vars and setting default values if missing
func FromEnv() (*Config, error) {
	host := loadEnvOrDefault("COMMENTS_API_HOST", "127.0.0.1")
	if host == "" {
		return nil, fmt.Errorf("host is missing")
	}

	port := loadEnvOrDefault("COMMENTS_API_PORT", "14000")
	if port == "" {
		return nil, fmt.Errorf("port is missing")
	}

	return &Config{
		Host: host,
		Port: port,
	}, nil
}
