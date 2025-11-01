package config

import (
    "context"
    "fmt"
    "os"
    "strconv"
    "sync"

    "github.com/jackc/pgx/v4/pgxpool"
    "github.com/joho/godotenv"
    "github.com/neo4j/neo4j-go-driver/v5/neo4j"

    "Loop_backend/internal/middleware"
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

    // Runtime connections
    PostgresPool    *pgxpool.Pool
    Neo4jDriver     neo4j.DriverWithContext
    AuthMiddleware  *middleware.AuthMiddleware
}

// ProviderType represents the type of AI provider
type ProviderType string

const (
	ProviderOllama  ProviderType = "ollama"
	ProviderOpenAI  ProviderType = "openai-compatible"
)

type AIConfig struct {
	// Provider Selection
	Provider ProviderType

	// // Ollama Configuration  
	// OllamaURL        string
	// OllamaModelName  string
	// OllamaEmbedModel string

	// OpenAI-compatible Configuration
	APIKey         string
	APIURL         string
	ModelName      string
	// EmbeddingModel string

	// Python Embedding Service Configuration
	PythonEmbeddingURL string
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
func LoadConfig(ctx context.Context) (*Config, error) {
    var loadErr error
    once.Do(func() {
        if err := godotenv.Load(); err != nil {
            fmt.Println("Warning: Error loading .env file")
        }

        // Initialize base config
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

            // // Ollama Configuration
            // OllamaURL:        getEnvValue("OLLAMA_URL", ""),
            // OllamaModelName:  getEnvValue("OLLAMA_MODEL_NAME", ""),
            // OllamaEmbedModel: getEnvValue("OLLAMA_EMBEDDING_MODEL", ""),

            // OpenAI-compatible Configuration
            APIKey:         getEnvValue("API_KEY", ""),
            APIURL:         getEnvValue("API_URL", ""),
            ModelName:      getEnvValue("MODEL_NAME", ""),
            // EmbeddingModel: getEnvValue("EMBEDDING_MODEL", ""),

            // Python Embedding Service Configuration
            PythonEmbeddingURL: getEnvValue("PYTHON_EMBEDDING_URL", "http://localhost:5000"),
        },
    }

        // Validate required environment variables
        if config.RelationalDatabaseConfig.ConnectionString == "" {
            loadErr = fmt.Errorf("DB_CONN_STRING environment variable is required")
            return
        }

        // Initialize connections
        var (
            postgresPool *pgxpool.Pool
            driver neo4j.DriverWithContext
            err error
        )

        // Initialize Postgres connection
        postgresPool, err = pgxpool.Connect(ctx, config.RelationalDatabaseConfig.ConnectionString)
        if err != nil {
            loadErr = fmt.Errorf("unable to connect to postgres: %v", err)
            return
        }
        config.PostgresPool = postgresPool

        // Initialize Neo4j connection
        driver, err = neo4j.NewDriverWithContext(
            config.Neo4jConfig.URI,
            neo4j.BasicAuth(config.Neo4jConfig.Username, config.Neo4jConfig.Password, ""),
        )
        if err != nil {
            loadErr = fmt.Errorf("unable to create neo4j driver: %v", err)
            return
        }
        config.Neo4jDriver = driver

        // Initialize auth middleware
        config.AuthMiddleware = middleware.NewAuthMiddleware(config.JWTConfig.Secret)
    })

    if loadErr != nil {
        return nil, loadErr
    }
    return &config, nil
}
