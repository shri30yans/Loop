package config

import (
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/joho/godotenv"
)

// supportedType is a constraint for types supported by getEnvValue
type supportedType interface {
	string | int
}

// getEnvValue returns environment variable value with proper type conversion
func getEnvValue[T supportedType](key string, defaultValue T) T {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}

	var result T
	switch any(defaultValue).(type) {
	case int:
		if intValue, err := strconv.Atoi(value); err == nil {
			result = any(intValue).(T)
		} else {
			result = defaultValue
		}
	case string:
		result = any(value).(T)
	default:
		result = defaultValue
	}
	return result
}

type Config struct {
	ServerConfig             ServerConfig
	RelationalDatabaseConfig RelationalDatabaseConfig
	Neo4jConfig              Neo4jConfig
	JWTConfig                JWTConfig
	AIConfig                 AIConfig
}

type AIConfig struct {
	OllamaURL                   string
	OllamaKey                   string
	OllamaModelName             string
	OllamaEmbedModel            string
	OllamaMaxAsync              int
	OllamaMaxTokenSize          int
	OllamaEmbeddingDim          int
	OllamaEmbeddingMaxTokenSize int
	OllamaNumCtx                int
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
				Port: getEnvValue("SERVER_PORT", 8080),
				Host: getEnvValue("SERVER_HOST", "localhost"),
			},
			RelationalDatabaseConfig: RelationalDatabaseConfig{
				ConnectionString: getEnvValue("DB_CONN_STRING", ""),
			},
			Neo4jConfig: Neo4jConfig{
				URI:      getEnvValue("NEO4J_URI", ""),
				Username: getEnvValue("NEO4J_USERNAME", "neo4j"),
				Password: getEnvValue("NEO4J_PASSWORD", ""),
			},
			JWTConfig: JWTConfig{
				Secret: getEnvValue("JWT_SECRET",""),
			},
			AIConfig: AIConfig{
				OllamaURL:                   getEnvValue("OLLAMA_URL", ""),
				OllamaKey:                   getEnvValue("OLLAMA_API_KEY", ""),
				OllamaModelName:             getEnvValue("OLLAMA_MODEL_NAME", "qwen2.5:3b"),
				OllamaEmbedModel:            getEnvValue("OLLAMA_EMBEDDING_MODEL", "nomic-embed-text"),
				OllamaMaxAsync:              getEnvValue("OLLAMA_MAX_ASYNC", 4),
				OllamaMaxTokenSize:          getEnvValue("OLLAMA_MAX_TOKEN_SIZE", 32768),
				OllamaEmbeddingDim:          getEnvValue("OLLAMA_EMBEDDING_DIM", 768),
				OllamaEmbeddingMaxTokenSize: getEnvValue("OLLAMA_EMBEDDING_MAX_TOKEN_SIZE", 8192),
				OllamaNumCtx:                getEnvValue("OLLAMA_NUM_CTX", 32768),
			},
		}

		if config.RelationalDatabaseConfig.ConnectionString == "" {
			err = fmt.Errorf("DB_CONN_STRING environment variable is required")
		}
	})

	return &config, nil
}
