package ollama

import (
    "Loop_backend/internal/ai/providers/utils"
)

// Provider implements the interfaces.Provider interface for Ollama
type Provider struct {
    httpClient     *utils.HTTPClient
    ChatModel      string
    EmbeddingModel string
}

// NewProvider creates a new Ollama provider instance
func NewProvider(apiBaseURL string, chatModel string, embeddingModel string) *Provider {
    return &Provider{
        httpClient:     utils.NewHTTPClient(apiBaseURL, ""), // Ollama doesn't use auth
        ChatModel:      chatModel,
        EmbeddingModel: embeddingModel,
    }
}
