package openai

import (
    "Loop_backend/internal/ai/providers/utils"
)

// Provider implements the interfaces.Provider interface for OpenAI-compatible APIs
type Provider struct {
    httpClient     *utils.HTTPClient
    ChatModel      string
    EmbeddingModel string
}

// NewProvider creates a new OpenAI provider instance
func NewProvider(apiBaseURL string, apiKey string, chatModel string, embeddingModel string) *Provider {
    return &Provider{
        httpClient:     utils.NewHTTPClient(apiBaseURL, apiKey),
        ChatModel:      chatModel,
        EmbeddingModel: embeddingModel,
    }
}
