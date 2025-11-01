package openai

import (
    "fmt"
    "Loop_backend/internal/ai/interfaces"
)

// apiChatRequest represents the request structure for chat completions
type apiChatRequest struct {
    Model       string             `json:"model"`
    Messages    []interfaces.Message `json:"messages"`
    Temperature float64            `json:"temperature,omitempty"`
}

// apiChatResponse represents the raw API response structure
type apiChatResponse struct {
    ID      string `json:"id"`
    Object  string `json:"object"`
    Created int64  `json:"created"`
    Model   string `json:"model"`
    Choices []struct {
        Index        int               `json:"index"`
        Message      interfaces.Message `json:"message"`
        FinishReason string           `json:"finish_reason"`
    } `json:"choices"`
    Usage struct {
        PromptTokens     int `json:"prompt_tokens"`
        CompletionTokens int `json:"completion_tokens"`
        TotalTokens      int `json:"total_tokens"`
    } `json:"usage"`
}

// Chat implements the interfaces.Provider Chat method
func (p *Provider) Chat(messages []interfaces.Message) (*interfaces.ChatResponse, error) {
    request := apiChatRequest{
        Model:       p.ChatModel,
        Messages:    messages,
        Temperature: 0.7,
    }

    var apiResp apiChatResponse
    err := p.httpClient.SendRequest("POST", "/v1/chat/completions", request, &apiResp, true)
    if err != nil {
        return nil, err
    }

    if len(apiResp.Choices) == 0 {
        return nil, fmt.Errorf("no response content available")
    }

    return &interfaces.ChatResponse{
        Content: apiResp.Choices[0].Message.Content,
    }, nil
}
