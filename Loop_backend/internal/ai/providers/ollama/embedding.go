package ollama

import (
    "fmt"
)

// GenerateEmbedding implements embedding generation
func (p *Provider) GenerateEmbedding(text string) ([]float64, error) {
    request := map[string]interface{}{
        "model": p.EmbeddingModel,
        "prompt": text,
    }

    var response struct {
        Embedding []float64 `json:"embedding"`
    }

    err := p.httpClient.SendRequest("POST", "/api/embeddings", request, &response, false)
    if err != nil {
        return nil, err
    }

    if len(response.Embedding) == 0 {
        return nil, fmt.Errorf("no embedding data available")
    }

    return response.Embedding, nil
}
