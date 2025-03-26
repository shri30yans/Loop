package utils

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
)

// HTTPClient provides common HTTP functionality for AI providers
type HTTPClient struct {
    BaseURL string
    APIKey  string
}

// NewHTTPClient creates a new HTTP client with the given base URL and optional API key
func NewHTTPClient(baseURL string, apiKey string) *HTTPClient {
    return &HTTPClient{
        BaseURL: baseURL,
        APIKey:  apiKey,
    }
}

// SendRequest handles HTTP requests to AI provider APIs
func (c *HTTPClient) SendRequest(method, endpoint string, body interface{}, target interface{}, withAuth bool) error {
    url := fmt.Sprintf("%s%s", c.BaseURL, endpoint)

    jsonBody, err := json.Marshal(body)
    if err != nil {
        return fmt.Errorf("failed to marshal request body: %v", err)
    }

    req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonBody))
    if err != nil {
        return fmt.Errorf("failed to create request: %v", err)
    }

    req.Header.Set("Content-Type", "application/json")
    if withAuth && c.APIKey != "" {
        req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.APIKey))
    }

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return fmt.Errorf("request failed: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("API returned status %d", resp.StatusCode)
    }

    return json.NewDecoder(resp.Body).Decode(target)
}
