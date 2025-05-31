package config

import (
	"os"
)

// Config holds all configuration for the application
type Config struct {
	Environment string
	Port        string
	DatabaseURL string
	
	// Database specific configs
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
}

// New creates a new config instance with values from environment variables
func New() *Config {
	return &Config{
		Environment: getEnv("ENVIRONMENT", "development"),
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", ""),
		
		// Database configs
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "password"),
		DBName:     getEnv("DB_NAME", "minilab"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
	}
}

// getEnv gets an environment variable with a fallback value
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// GetDatabaseURL returns the complete database URL
func (c *Config) GetDatabaseURL() string {
	if c.DatabaseURL != "" {
		return c.DatabaseURL
	}
	
	return "postgres://" + c.DBUser + ":" + c.DBPassword + 
		   "@" + c.DBHost + ":" + c.DBPort + 
		   "/" + c.DBName + "?sslmode=" + c.DBSSLMode
} 