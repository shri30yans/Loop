package config

import (
	"fmt"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

func getEnvOrDefault(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}

type Config struct {
	ServerConfig   ServerConfig
	RelationalDatabaseConfig RelationalDatabaseConfig
	Neo4jConfig    Neo4jConfig
	JWTConfig      JWTConfig
}

type ServerConfig struct {
	Port int
	Host string
}

type RelationalDatabaseConfig struct {
	ConnectionString string
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
			ServerConfig: ServerConfig{
				Port: 8080,
				Host: getEnvOrDefault("SERVER_HOST", "localhost"),
			},
			RelationalDatabaseConfig: RelationalDatabaseConfig{
				ConnectionString: getEnvOrDefault("DB_CONN_STRING", ""),
			},
			Neo4jConfig: Neo4jConfig{
				URI:      getEnvOrDefault("NEO4J_URI", "neo4j://localhost:7687"),
				Username: getEnvOrDefault("NEO4J_USERNAME", "neo4j"),
				Password: os.Getenv("NEO4J_PASSWORD"),
			},
			JWTConfig: JWTConfig{
				Secret: os.Getenv("JWT_SECRET"),
			},
		}

		
		if config.RelationalDatabaseConfig.ConnectionString == "" {
			err = fmt.Errorf("DB_CONN_STRING environment variable is required")
		}
		
	})

	return &config, nil
}
