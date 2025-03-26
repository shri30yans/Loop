package providers

import (
    "fmt"
    "Loop_backend/config"
    "Loop_backend/internal/ai/interfaces"
    "Loop_backend/internal/ai/providers/ollama"
    "Loop_backend/internal/ai/providers/openai"
)

// NewProvider creates a new AI provider based on configuration
func NewProvider(cfg *config.AIConfig) (interfaces.Provider, error) {
    switch cfg.Provider {
    case config.ProviderOllama:
        if cfg.OllamaURL == "" {
            return nil, fmt.Errorf("ollama URL is required for Ollama provider")
        }
        return ollama.NewProvider(
            cfg.OllamaURL,
            cfg.OllamaModelName,
            cfg.OllamaEmbedModel,
        ), nil

    case config.ProviderOpenAI:
        if cfg.APIKey == "" {
            return nil, fmt.Errorf("API key is required for %s provider", cfg.Provider)
        }
        // All these providers use OpenAI-compatible API
        return openai.NewProvider(
            cfg.APIURL,
            cfg.APIKey,
            cfg.ModelName,
            cfg.EmbeddingModel,
        ), nil

    default:
        return nil, fmt.Errorf("unsupported provider type: %s", cfg.Provider)
    }
}
