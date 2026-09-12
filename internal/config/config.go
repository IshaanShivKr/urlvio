package config

import (
	"fmt"
	"os"
)

const (
	defaultPort    = "8080"
	defaultGinMode = "debug"
	defaultBaseURL = "http://localhost:8080"
)

type Config struct {
	BaseURL     string
	DatabaseURL string
	GinMode     string
	Port        string
}

func Load() (*Config, error) {
	cfg := &Config{
		BaseURL:     os.Getenv("BASE_URL"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
		GinMode:     os.Getenv("GIN_MODE"),
		Port:        os.Getenv("PORT"),
	}

	if cfg.BaseURL == "" {
		cfg.BaseURL = defaultBaseURL
	}

	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is required")
	}

	if cfg.GinMode == "" {
		cfg.GinMode = defaultGinMode
	}

	if cfg.Port == "" {
		cfg.Port = defaultPort
	}

	return cfg, nil
}
