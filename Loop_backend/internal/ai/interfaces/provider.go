package interfaces

// Provider defines the interface for LLM providers
type Provider interface {
    Chat(messages []Message) (*ChatResponse, error)
    GenerateEmbedding(input string) ([]float64, error)
}
