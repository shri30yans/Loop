package embedding

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// PythonEmbeddingClient implements EmbeddingClient using HTTP requests to Python server
type PythonEmbeddingClient struct {
	baseURL    string
	httpClient *http.Client
}

// PythonEmbeddingRequest represents the request body for embedding generation
type PythonEmbeddingRequest struct {
	Text  string   `json:"text,omitempty"`
	Texts []string `json:"texts,omitempty"`
}

// PythonEmbeddingResponse represents the response from the Python embedding service
type PythonEmbeddingResponse struct {
	Embedding        []float32   `json:"embedding,omitempty"`
	Embeddings       [][]float32 `json:"embeddings,omitempty"`
	Dimensions       int         `json:"dimensions"`
	ProcessingTimeMs float64     `json:"processing_time_ms"`
	Error            string      `json:"error,omitempty"`
}

// PythonHealthResponse represents the response from the health check endpoint
type PythonHealthResponse struct {
	Status     string `json:"status"`
	Model      string `json:"model"`
	Dimensions int    `json:"dimensions"`
	GoBackend  string `json:"go_backend"`
}

// NewPythonEmbeddingClient creates a new Python embedding client
func NewPythonEmbeddingClient(baseURL string) EmbeddingClient {
	return &PythonEmbeddingClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GenerateEmbedding generates an embedding for a single text
func (c *PythonEmbeddingClient) GenerateEmbedding(text string) ([]float32, error) {
	// Create request body
	reqBody, err := json.Marshal(PythonEmbeddingRequest{Text: text})
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
	var embeddingResp PythonEmbeddingResponse
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
func (c *PythonEmbeddingClient) GenerateBatchEmbeddings(texts []string) ([][]float32, error) {
	// Create request body
	reqBody, err := json.Marshal(PythonEmbeddingRequest{Texts: texts})
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
	var embeddingResp PythonEmbeddingResponse
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
func (c *PythonEmbeddingClient) HealthCheck() (bool, error) {
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
	var healthResp PythonHealthResponse
	if err := json.NewDecoder(resp.Body).Decode(&healthResp); err != nil {
		return false, fmt.Errorf("failed to decode response: %w", err)
	}

	return healthResp.Status == "healthy", nil
}
