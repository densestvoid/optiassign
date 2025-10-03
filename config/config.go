package config

import (
	"context"
	"time"

	"github.com/sethvargo/go-envconfig"
)

// Config holds all application configuration
type Config struct {
	// Server configuration
	Port string `env:"PORT,default=8080"`
	
	// Database configuration
	DatabaseURL string `env:"DATABASE_URL,required"`
	
	// Google OAuth configuration
	GoogleClientID     string `env:"GOOGLE_CLIENT_ID,required"`
	GoogleClientSecret string `env:"GOOGLE_CLIENT_SECRET,required"`
	GoogleRedirectURL  string `env:"GOOGLE_REDIRECT_URL,required"`
	
	// Session configuration
	SessionSecret string `env:"SESSION_SECRET,required"`
	SessionTTL    time.Duration `env:"SESSION_TTL,default=24h"`
	
	// Email configuration (for development)
	EmailLogToConsole bool `env:"EMAIL_LOG_TO_CONSOLE,default=true"`
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	var cfg Config
	if err := envconfig.Process(context.Background(), &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}