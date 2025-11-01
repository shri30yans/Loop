package search

import (
	"Loop_backend/internal/dto"
	"Loop_backend/internal/repositories/vector"
	"Loop_backend/internal/services/embedding"
	"context"
	"fmt"
	"github.com/google/uuid"
)

// VectorSearchService implements the VectorService interface
type VectorSearchService struct {
	vectorRepo      vector.VectorRepository
	embeddingClient embedding.EmbeddingClient
}

// NewVectorSearchService creates a new vector search service
func NewVectorSearchService(
	vectorRepo vector.VectorRepository,
	embeddingClient embedding.EmbeddingClient,
) *VectorSearchService {
	return &VectorSearchService{
		vectorRepo:      vectorRepo,
		embeddingClient: embeddingClient,
	}
}

// SearchProjects performs a vector-based search for projects
func (s *VectorSearchService) SearchProjects(ctx context.Context, query string) ([]dto.SearchToolsResult, error) {
	// Generate embedding for the query
	embedding, err := s.embeddingClient.GenerateEmbedding(query)
	if err != nil {
		return nil, fmt.Errorf("failed to generate embedding: %w", err)
	}

	// Search for similar projects
	projectIDs, err := s.vectorRepo.SearchSimilarProjects(embedding, 10) // Limit to 10 results
	if err != nil {
		return nil, fmt.Errorf("failed to search similar projects: %w", err)
	}

	// Convert project IDs to search results
	results := make([]dto.SearchToolsResult, 0, len(projectIDs))
	for i, projectID := range projectIDs {
		// Calculate a simple score based on position (higher is better)
		score := float64(len(projectIDs)-i) / float64(len(projectIDs))
		
		results = append(results, dto.SearchToolsResult{
			Score:     score,
			ProjectID: projectID.String(),
		})
	}

	return results, nil
}

// GetSimilarProjects finds projects similar to the given project
func (s *VectorSearchService) GetSimilarProjects(ctx context.Context, projectID string) ([]dto.SearchToolsResult, error) {
	// Parse project ID
	id, err := uuid.Parse(projectID)
	if err != nil {
		return nil, fmt.Errorf("invalid project ID: %w", err)
	}

	// Get the project's vector
	vector, err := s.vectorRepo.GetProjectVector(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get project vector: %w", err)
	}

	// Search for similar projects
	similarIDs, err := s.vectorRepo.SearchSimilarProjects(vector, 10) // Limit to 10 results
	if err != nil {
		return nil, fmt.Errorf("failed to search similar projects: %w", err)
	}

	// Convert project IDs to search results
	results := make([]dto.SearchToolsResult, 0, len(similarIDs))
	for i, similarID := range similarIDs {
		// Skip the original project
		if similarID == id {
			continue
		}

		// Calculate a simple score based on position (higher is better)
		score := float64(len(similarIDs)-i) / float64(len(similarIDs))
		
		results = append(results, dto.SearchToolsResult{
			Score:     score,
			ProjectID: similarID.String(),
		})
	}

	return results, nil
}
