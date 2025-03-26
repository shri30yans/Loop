package ollama

import (
    "fmt"
    "Loop_backend/internal/ai/interfaces"
)

// Chat implements the Provider Chat method
func (p *Provider) Chat(messages []interfaces.Message) (*interfaces.ChatResponse, error) {
    // Convert messages to Ollama format
    prompt := ""
    for _, msg := range messages {
        if msg.Role == "system" {
            prompt += "System: " + msg.Content + "\n"
        } else if msg.Role == "assistant" {
            prompt += "Assistant: " + msg.Content + "\n"
        } else {
            prompt += "Human: " + msg.Content + "\n"
        }
    }

    request := map[string]interface{}{
        "model":       p.ChatModel,
        "prompt":      prompt,
        "stream":      false,
        "temperature": 0.7,
    }

    var response struct {
        Response string `json:"response"`
    }

    err := p.httpClient.SendRequest("POST", "/api/generate", request, &response, false)
    if err != nil {
        return nil, err
    }

    if response.Response == "" {
        return nil, fmt.Errorf("no response content available")
    }

    return &interfaces.ChatResponse{
        Content: response.Response,
    }, nil
}
