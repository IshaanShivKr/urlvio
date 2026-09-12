package config

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	BaseURL     string
	DatabaseURL string
	GinMode     string
	Port        string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		BaseURL:     os.Getenv("BASE_URL"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
		GinMode:     os.Getenv("GIN_MODE"),
		Port:        os.Getenv("PORT"),
	}

	if err := validate(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func validate(cfg *Config) error {
	if cfg.BaseURL == "" {
		return fmt.Errorf("BASE_URL environment variable is required")
	}

	if err := validateBaseURL(cfg.BaseURL); err != nil {
		return err
	}

	if cfg.DatabaseURL == "" {
		return fmt.Errorf("DATABASE_URL environment variable is required")
	}

	if cfg.GinMode == "" {
		return fmt.Errorf("GIN_MODE environment variable is required")
	}

	switch cfg.GinMode {
	case "debug", "release", "test":
	default:
		return fmt.Errorf("GIN_MODE must be one of: debug, release, test")
	}

	if cfg.Port == "" {
		return fmt.Errorf("PORT environment variable is required")
	}

	port, err := strconv.Atoi(cfg.Port)
	if err != nil || port < 1 || port > 65535 {
		return fmt.Errorf("PORT must be a valid port between 1 and 65535")
	}

	return nil
}

func validateBaseURL(rawURL string) error {
	u, err := url.Parse(rawURL)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return fmt.Errorf("BASE_URL must be a valid absolute URL")
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("BASE_URL must use http or https")
	}

	if strings.Contains(u.Path, "//") {
		return fmt.Errorf("BASE_URL must not contain an invalid path")
	}

	return nil
}
