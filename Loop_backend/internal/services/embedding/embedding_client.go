package embedding

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// EmbeddingClient defines the interface for generating embeddings
type EmbeddingClient interface {
	GenerateEmbedding(text string) ([]float32, error)
	GenerateBatchEmbeddings(texts []string) ([][]float32, error)
	HealthCheck() (bool, error)
}

// HTTPEmbeddingClient implements EmbeddingClient using HTTP requests
type HTTPEmbeddingClient struct {
	baseURL    string
	httpClient *http.Client
}

// EmbeddingRequest represents the request body for embedding generation
type EmbeddingRequest struct {
	Text  string   `json:"text,omitempty"`
	Texts []string `json:"texts,omitempty"`
}

// EmbeddingResponse represents the response from the embedding service
type EmbeddingResponse struct {
	Embedding  []float32   `json:"embedding,omitempty"`
	Embeddings [][]float32 `json:"embeddings,omitempty"`
	Dimensions int         `json:"dimensions"`
	Error      string      `json:"error,omitempty"`
}

// HealthResponse represents the response from the health check endpoint
type HealthResponse struct {
	Status     string `json:"status"`
	Model      string `json:"model"`
	Dimensions int    `json:"dimensions"`
}

// NewEmbeddingClient creates a new HTTP embedding client
func NewEmbeddingClient(baseURL string) EmbeddingClient {
	return &HTTPEmbeddingClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GenerateEmbedding generates an embedding for a single text
func (c *HTTPEmbeddingClient) GenerateEmbedding(text string) ([]float32, error) {
	// Create request body
	reqBody, err := json.Marshal(EmbeddingRequest{Text: text})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", c.baseURL+"/embed", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("embedding service returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var embeddingResp EmbeddingResponse
	if err := json.Unmarshal(body, &embeddingResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Check for error in response
	if embeddingResp.Error != "" {
		return nil, fmt.Errorf("embedding service error: %s", embeddingResp.Error)
	}

	return embeddingResp.Embedding, nil
}

// GenerateBatchEmbeddings generates embeddings for multiple texts
func (c *HTTPEmbeddingClient) GenerateBatchEmbeddings(texts []string) ([][]float32, error) {
	// Create request body
	reqBody, err := json.Marshal(EmbeddingRequest{Texts: texts})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", c.baseURL+"/embed-batch", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("embedding service returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var embeddingResp EmbeddingResponse
	if err := json.Unmarshal(body, &embeddingResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Check for error in response
	if embeddingResp.Error != "" {
		return nil, fmt.Errorf("embedding service error: %s", embeddingResp.Error)
	}

	return embeddingResp.Embeddings, nil
}

// HealthCheck checks if the embedding service is healthy
func (c *HTTPEmbeddingClient) HealthCheck() (bool, error) {
	// Create HTTP request
	req, err := http.NewRequest("GET", c.baseURL+"/health", nil)
	if err != nil {
		return false, fmt.Errorf("failed to create request: %w", err)
	}

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("embedding service returned status %d", resp.StatusCode)
	}

	// Parse response
	var healthResp HealthResponse
	if err := json.NewDecoder(resp.Body).Decode(&healthResp); err != nil {
		return false, fmt.Errorf("failed to decode response: %w", err)
	}

	return healthResp.Status == "healthy", nil
}
