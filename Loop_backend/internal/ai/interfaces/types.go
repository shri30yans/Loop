package interfaces

// Message represents a chat message
type Message struct {
    Role    string `json:"role"`
    Content string `json:"content"`
}

// ChatResponse represents a response from the LLM
type ChatResponse struct {
    Content string `json:"content"`
}

// ProviderType represents the type of AI provider
type ProviderType string

const (
    OpenAI ProviderType = "openai"
    Ollama ProviderType = "ollama"
)

// Config represents configuration for an AI provider
type Config struct {
    Type           ProviderType `json:"type"`
    APIURL         string       `json:"api_url"`
    APIKey         string       `json:"api_key,omitempty"`
    ChatModel      string       `json:"chat_model"`
    EmbeddingModel string       `json:"embedding_model"`
}
