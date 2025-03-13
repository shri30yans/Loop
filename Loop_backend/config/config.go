package config

import (
	"fmt"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

type ServerConfig struct {
	Port int
	Host string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type JWTConfig struct {
	Secret string
}

var (
	config Config
	once   sync.Once
)

/*
LoadConfig loads configuration from environment variables.
It uses the singleton pattern to ensure the configuration is loaded only once.
*/
func LoadConfig() (*Config, error) {
	var err error
	once.Do(func() {
		err = godotenv.Load()
		if err != nil {
			fmt.Println("Warning: Error loading .env file")
		}

		config = Config{
			Server: ServerConfig{
				Port: 8080,
				Host: getEnvOrDefault("SERVER_HOST", "localhost"),
			},
			Database: DatabaseConfig{
				Host:     getEnvOrDefault("DB_HOST", "localhost"),
				Port:     getEnvOrDefault("DB_PORT", "5432"),
				User:     getEnvOrDefault("DB_USER", "postgres"),
				Password: os.Getenv("DB_PASSWORD"),
				Name:     getEnvOrDefault("DB_NAME", "loop"),
			},
			JWT: JWTConfig{
				Secret: os.Getenv("JWT_SECRET"),
			},
		}

		// Validate required configurations
		if config.Database.Password == "" {
			err = fmt.Errorf("DB_PASSWORD environment variable is required")
		}
		if config.JWT.Secret == "" {
			err = fmt.Errorf("JWT_SECRET environment variable is required")
		}
	})

	if err != nil {
		return nil, err
	}

	return &config, nil
}

// GetConfig returns the current configuration
func GetConfig() *Config {
	return &config
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
