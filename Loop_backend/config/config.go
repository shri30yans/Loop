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

// ProviderType represents the type of AI provider
type ProviderType string

const (
	ProviderOllama ProviderType = "ollama"
	ProviderOpenAI ProviderType = "openai-compatible"
)

type AIConfig struct {
	// Provider Selection
	Provider ProviderType

	// Ollama Configuration
	OllamaURL        string
	OllamaModelName  string
	OllamaEmbedModel string

	// OpenAI-compatible Configuration
	APIKey         string
	APIURL         string
	ModelName      string
	EmbeddingModel string
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
				Secret: getEnvValue("JWT_SECRET", ""),
			},
			AIConfig: AIConfig{
				Provider: ProviderType(getEnvValue("AI_PROVIDER", string(ProviderOpenAI))),

				// Ollama Configuration
				OllamaURL:        getEnvValue("OLLAMA_URL", ""),
				OllamaModelName:  getEnvValue("OLLAMA_MODEL_NAME", ""),
				OllamaEmbedModel: getEnvValue("OLLAMA_EMBEDDING_MODEL", ""),

				// OpenAI-compatible Configuration
				APIKey:         getEnvValue("API_KEY", ""),
				APIURL:         getEnvValue("API_URL", ""),
				ModelName:      getEnvValue("MODEL_NAME", ""),
				EmbeddingModel: getEnvValue("EMBEDDING_MODEL", ""),
			},
		}

		if config.RelationalDatabaseConfig.ConnectionString == "" {
			err = fmt.Errorf("DB_CONN_STRING environment variable is required")
		}
	})

	return &config, nil
}
