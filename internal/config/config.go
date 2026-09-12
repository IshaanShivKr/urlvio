package config

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	BaseURL     string
	DatabaseURL string
	GinMode     string
	Port        string
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		slog.Warn("could not load .env file", "error", err)
	}

	cfg := &Config{
		BaseURL:     os.Getenv("BASE_URL"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
		GinMode:     os.Getenv("GIN_MODE"),
		Port:        os.Getenv("PORT"),
	}

	if cfg.BaseURL == "" {
		return nil, fmt.Errorf("BASE_URL environment variable is required")
	}

	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is required")
	}

	if cfg.GinMode == "" {
		return nil, fmt.Errorf("GIN_MODE environment variable is required")
	}

	if cfg.Port == "" {
		return nil, fmt.Errorf("PORT environment variable is required")
	}

	return cfg, nil
}
