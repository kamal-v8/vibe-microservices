// Package config handles environment-based configuration for the user service.
// Follows the 12-factor app methodology — all config comes from environment variables.
package config

import (
	"os"
)

// Config holds all configuration values for the user service.
type Config struct {
	Port        string // HTTP server port
	DatabaseURL string // PostgreSQL connection string
	LogLevel    string // Logging verbosity (info, debug, warn, error)
}

// Load reads configuration from environment variables, falling back to sensible
// defaults when a variable is not set. DATABASE_URL has no default because
// running without a database is a fatal misconfiguration — the caller should
// check for an empty value.
func Load() *Config {
	return &Config{
		Port:        getEnv("PORT", "8081"),
		DatabaseURL: getEnv("DATABASE_URL", ""),
		LogLevel:    getEnv("LOG_LEVEL", "info"),
	}
}

// getEnv retrieves an environment variable or returns the provided fallback.
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
