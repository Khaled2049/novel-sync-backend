// File: internal/config/config.go
package config

import (
	"os"
	"strconv"
	"time"
	// If using Viper or godotenv, import them here
	// "github.com/spf13/viper"
	// "github.com/joho/godotenv"
)

type JWTConfig struct {
	SecretKey string        `mapstructure:"secretKey"` // Should be loaded securely!
	TTL       time.Duration `mapstructure:"ttl"`       // Token Time-To-Live
}

// Config holds all configuration for the application.
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Firebase FirebaseConfig `mapstructure:"firebase"`
	Database DatabaseConfig `mapstructure:"database"` // Added Database config
	// RabbitMQ RabbitMQConfig `mapstructure:"rabbitmq"`
	JWT      JWTConfig      `mapstructure:"jwt"`
}

// ServerConfig holds HTTP server specific configuration.
type ServerConfig struct {
	Port         string        `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"readTimeout"`
	WriteTimeout time.Duration `mapstructure:"writeTimeout"`
	IdleTimeout  time.Duration `mapstructure:"idleTimeout"`
}

// FirebaseConfig holds Firebase specific configuration.
type FirebaseConfig struct {
	ServiceAccountKeyPath string `mapstructure:"serviceAccountKeyPath"`
}

// DatabaseConfig holds database specific configuration.
type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbName"`
	SSLMode  string `mapstructure:"sslMode"`
}

// LoadConfig reads configuration from file or environment variables.
// --- Updated Placeholder LoadConfig ---
func LoadConfig() (*Config, error) {
	// In a real app, use godotenv.Load() here if using .env directly with Go
	// Or setup Viper to read env vars / config files

	// Example reading directly from Env Vars set by .env for local Go run
	dbPortStr := getEnv("APP_DB_PORT", "5433") // Default to docker host port
	dbPort, err := strconv.Atoi(dbPortStr)
	if err != nil {
		dbPort = 5433 // Fallback to default port on error
	}

	jwtSecret := getEnv("JWT_SECRET_KEY", "default-super-secret-key") // !! CHANGE THIS & LOAD SECURELY !!
	jwtTTLStr := getEnv("JWT_TTL_MINUTES", "60")                      // Default to 60 minutes
	jwtTTLMinutes, err := strconv.Atoi(jwtTTLStr)
	if err != nil {
		jwtTTLMinutes = 60 // Fallback
	}

	return &Config{
		Server: ServerConfig{
			Port:         getEnv("APP_SERVER_PORT", "8000"),
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
		Firebase: FirebaseConfig{
			ServiceAccountKeyPath: getEnv("FIREBASE_SERVICE_ACCOUNT_KEY_PATH", ""),
		},
		Database: DatabaseConfig{
			Host:     getEnv("APP_DB_HOST", "localhost"),
			Port:     dbPort,
			User:     getEnv("APP_DB_USER", "user"), // Match .env default user
			Password: getEnv("APP_DB_PASSWORD", "password"), // Match .env default password
			DBName:   getEnv("APP_DB_NAME", "noveldb"),   // Match .env default dbname
			SSLMode:  getEnv("APP_DB_SSL_MODE", "disable"), // Default to disable for local docker
		},
		JWT: JWTConfig{
			SecretKey: jwtSecret,
			TTL:       time.Duration(jwtTTLMinutes) * time.Minute,
		},
		// Initialize other configs
	}, nil
}

// Helper function to get env var or default
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}