package openai

import (
    "fmt"
)

// EmbeddingRequest represents the request structure for embeddings
type EmbeddingRequest struct {
    Model string `json:"model"`
    Input string `json:"input"`
}

// EmbeddingResponse represents the response structure from embeddings
type EmbeddingResponse struct {
    Object string `json:"object"`
    Data   []struct {
        Object    string    `json:"object"`
        Embedding []float64 `json:"embedding"`
        Index     int       `json:"index"`
    } `json:"data"`
    Model string `json:"model"`
    Usage struct {
        PromptTokens int `json:"prompt_tokens"`
        TotalTokens  int `json:"total_tokens"`
    } `json:"usage"`
}

// GenerateEmbedding implements embedding generation
func (p *Provider) GenerateEmbedding(text string) ([]float64, error) {
    request := EmbeddingRequest{
        Model: p.EmbeddingModel,
        Input: text,
    }

    var response EmbeddingResponse
    err := p.httpClient.SendRequest("POST", "/v1/embeddings", request, &response, true)
    if err != nil {
        return nil, err
    }

    if len(response.Data) == 0 {
        return nil, fmt.Errorf("no embedding data available")
    }

    return response.Data[0].Embedding, nil
}
