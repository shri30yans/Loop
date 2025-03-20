package config

import (
	"fmt"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
<<<<<<< HEAD
	ServerConfig   ServerConfig
	RelationalDatabaseConfig RelationalDatabaseConfig
	Neo4jConfig    Neo4jConfig
	JWTConfig      JWTConfig
=======
	Server   ServerConfig
	Database DatabaseConfig
	Neo4j    Neo4jConfig
	JWT      JWTConfig
>>>>>>> 4a2f436bed91636c5c2e3782993f5ab211ecfca7
}

type ServerConfig struct {
	Port int
	Host string
}

<<<<<<< HEAD
type RelationalDatabaseConfig struct {
	ConnectionString string
=======
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
>>>>>>> 4a2f436bed91636c5c2e3782993f5ab211ecfca7
}

type Neo4jConfig struct {
	URI      string
	Username string
	Password string
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
<<<<<<< HEAD
			ServerConfig: ServerConfig{
				Port: 8080,
				Host: getEnvOrDefault("SERVER_HOST", "localhost"),
			},
			RelationalDatabaseConfig: RelationalDatabaseConfig{
				ConnectionString: getEnvOrDefault("DB_CONN_STRING", ""),
			},
			Neo4jConfig: Neo4jConfig{
=======
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
			Neo4j: Neo4jConfig{
>>>>>>> 4a2f436bed91636c5c2e3782993f5ab211ecfca7
				URI:      getEnvOrDefault("NEO4J_URI", "neo4j://localhost:7687"),
				Username: getEnvOrDefault("NEO4J_USERNAME", "neo4j"),
				Password: os.Getenv("NEO4J_PASSWORD"),
			},
<<<<<<< HEAD
			JWTConfig: JWTConfig{
=======
			JWT: JWTConfig{
>>>>>>> 4a2f436bed91636c5c2e3782993f5ab211ecfca7
				Secret: os.Getenv("JWT_SECRET"),
			},
		}

		// Validate required configurations
<<<<<<< HEAD
		if config.RelationalDatabaseConfig.ConnectionString == "" {
			err = fmt.Errorf("DB_CONN_STRING environment variable is required")
		}
		
=======
		if config.Database.Password == "" {
			err = fmt.Errorf("DB_PASSWORD environment variable is required")
		}
		if config.JWT.Secret == "" {
			err = fmt.Errorf("JWT_SECRET environment variable is required")
		}
		if config.Neo4j.Password == "" {
			err = fmt.Errorf("NEO4J_PASSWORD environment variable is required")
		}
>>>>>>> 4a2f436bed91636c5c2e3782993f5ab211ecfca7
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
