package config

import (
	"fmt"
	"os"
	"strings"
)

// Environment represents the application environment
type Environment string

const (
	Development Environment = "development"
	Staging     Environment = "staging"
	Production  Environment = "production"
)

// GetEnvironment returns the current environment
func GetEnvironment() Environment {
	env := os.Getenv("APP_ENV")
	switch strings.ToLower(env) {
	case "production":
		return Production
	case "staging":
		return Staging
	default:
		return Development
	}
}

// IsDevelopment returns true if running in development
func IsDevelopment() bool {
	return GetEnvironment() == Development
}

// IsProduction returns true if running in production
func IsProduction() bool {
	return GetEnvironment() == Production
}

// GetLogLevel returns the log level based on environment
func GetLogLevel() string {
	if IsProduction() {
		return "info"
	}
	return "debug"
}

// GetDatabaseMaxConnections returns max database connections based on environment
func GetDatabaseMaxConnections() int {
	if IsProduction() {
		return 25
	}
	return 10
}

// GetTaskManagerWorkers returns number of background workers based on environment
func GetTaskManagerWorkers() int {
	if IsProduction() {
		return 5
	}
	return 2
}

// GetRateLimit returns rate limit based on environment
func GetRateLimit() int {
	if IsProduction() {
		return 100 // 100 requests per minute
	}
	return 1000 // 1000 requests per minute for development
}

// GetSessionTimeout returns session timeout based on environment
func GetSessionTimeout() int {
	if IsProduction() {
		return 3600 // 1 hour
	}
	return 86400 // 24 hours for development
}

// GetEmailLogToConsole returns whether to log emails to console
func GetEmailLogToConsole() bool {
	if IsProduction() {
		return false // In production, use real email service
	}
	return true // In development, log to console
}

// GetBaseURL returns the base URL based on environment
func GetBaseURL() string {
	if IsProduction() {
		return os.Getenv("BASE_URL")
	}
	return "http://localhost:8080"
}

// GetDatabaseURL returns the database URL with environment-specific settings
func GetDatabaseURL() string {
	baseURL := os.Getenv("DATABASE_URL")
	if baseURL == "" {
		// Construct from individual components
		host := os.Getenv("DB_HOST")
		port := os.Getenv("DB_PORT")
		user := os.Getenv("DB_USER")
		password := os.Getenv("DB_PASSWORD")
		dbname := os.Getenv("DB_NAME")
		sslmode := "require"
		if IsDevelopment() {
			sslmode = "disable"
		}
		
		baseURL = "postgres://" + user + ":" + password + "@" + host + ":" + port + "/" + dbname + "?sslmode=" + sslmode
	}
	return baseURL
}

// GetGoogleOAuthConfig returns Google OAuth configuration
func GetGoogleOAuthConfig() (clientID, clientSecret, redirectURL string) {
	clientID = os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret = os.Getenv("GOOGLE_CLIENT_SECRET")
	redirectURL = GetBaseURL() + "/auth/google/callback"
	return
}

// GetSessionConfig returns session configuration
func GetSessionConfig() (secret string, timeout int) {
	secret = os.Getenv("SESSION_SECRET")
	if secret == "" {
		secret = "default-secret-key-change-in-production"
	}
	timeout = GetSessionTimeout()
	return
}

// ValidateProductionConfig validates production configuration
func ValidateProductionConfig() error {
	if !IsProduction() {
		return nil
	}
	
	requiredVars := []string{
		"DB_HOST", "DB_PORT", "DB_NAME", "DB_USER", "DB_PASSWORD",
		"GOOGLE_CLIENT_ID", "GOOGLE_CLIENT_SECRET", "SESSION_SECRET",
	}
	
	for _, varName := range requiredVars {
		if os.Getenv(varName) == "" {
			return fmt.Errorf("required environment variable %s is not set", varName)
		}
	}
	
	return nil
}

// GetPort returns the port to listen on
func GetPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return port
}

// GetHTTPSEnabled returns whether HTTPS is enabled
func GetHTTPSEnabled() bool {
	return os.Getenv("HTTPS_ENABLED") == "true"
}

// GetTrustedProxies returns trusted proxy IPs
func GetTrustedProxies() []string {
	proxies := os.Getenv("TRUSTED_PROXIES")
	if proxies == "" {
		return []string{}
	}
	return strings.Split(proxies, ",")
}