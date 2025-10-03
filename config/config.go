package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds all application configuration
type Config struct {
	// Server configuration
	Port string
	
	// Database configuration
	DatabaseURL string
	
	// Google OAuth configuration
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string
	
	// Session configuration
	SessionSecret string
	SessionTTL    time.Duration
	
	// Email configuration (for development)
	EmailLogToConsole bool
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	cfg := &Config{
		Port:               getEnv("PORT", "8080"),
		DatabaseURL:        getEnv("DATABASE_URL", ""),
		GoogleClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
		GoogleRedirectURL:  getEnv("GOOGLE_REDIRECT_URL", ""),
		SessionSecret:      getEnv("SESSION_SECRET", ""),
		SessionTTL:        24 * time.Hour,
		EmailLogToConsole:  getEnvBool("EMAIL_LOG_TO_CONSOLE", true),
	}
	
	// Validate required fields
	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}
	if cfg.GoogleClientID == "" {
		return nil, fmt.Errorf("GOOGLE_CLIENT_ID is required")
	}
	if cfg.GoogleClientSecret == "" {
		return nil, fmt.Errorf("GOOGLE_CLIENT_SECRET is required")
	}
	if cfg.GoogleRedirectURL == "" {
		return nil, fmt.Errorf("GOOGLE_REDIRECT_URL is required")
	}
	if cfg.SessionSecret == "" {
		return nil, fmt.Errorf("SESSION_SECRET is required")
	}
	
	return cfg, nil
}

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvBool gets a boolean environment variable with a default value
func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseBool(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}