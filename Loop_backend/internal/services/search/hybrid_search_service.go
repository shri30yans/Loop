package search

import (
	"Loop_backend/internal/dto"
	"Loop_backend/internal/services/graph"
	"context"
	"fmt"
	"sort"
)

// HybridSearchService combines vector and graph-based search strategies
type HybridSearchService struct {
	graphService  graph.GraphService
	vectorService VectorService
}

// NewHybridSearchService creates a new hybrid search service
func NewHybridSearchService(
	graphService graph.GraphService,
	vectorService VectorService,
) *HybridSearchService {
	return &HybridSearchService{
		graphService:  graphService,
		vectorService: vectorService,
	}
}

// VectorSearch performs a vector-based search
func (s *HybridSearchService) VectorSearch(ctx context.Context, query string) ([]dto.SearchToolsResult, error) {
	if s.vectorService == nil {
		return nil, fmt.Errorf("vector service not initialized")
	}
	return s.vectorService.SearchProjects(ctx, query)
}

// GraphSearch performs a graph-based search
func (s *HybridSearchService) GraphSearch(ctx context.Context, query string) ([]dto.SearchToolsResult, error) {
	if s.graphService == nil {
		// Fall back to vector search if graph service is not available
		return s.VectorSearch(ctx, query)
	}
	return s.graphService.SearchProjects(ctx, query)
}

// HybridSearch performs a hybrid search across all strategies
func (s *HybridSearchService) HybridSearch(ctx context.Context, query string) ([]dto.SearchToolsResult, error) {
	if s.vectorService == nil {
		return nil, fmt.Errorf("vector service not initialized")
	}

	// When graph service is not available, fall back to pure vector search
	if s.graphService == nil {
		return s.VectorSearch(ctx, query)
	}

	// Get base results from vector search
	baseResults, err := s.vectorService.SearchProjects(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("vector search failed: %w", err)
	}

	// Convert base results to project IDs
	var projectIDs []string
	for _, result := range baseResults {
		projectIDs = append(projectIDs, result.ProjectID)
	}

	// Get metadata for projects using graph service
	metadata, err := s.graphService.GetProjectsMetadata(ctx, projectIDs)
	if err != nil {
		// If metadata retrieval fails, return vector results without enrichment
		return baseResults, nil
	}

	// Enhance results with metadata
	var results []dto.SearchToolsResult
	for _, result := range baseResults {
		if m, ok := metadata[result.ProjectID]; ok {
			results = append(results, dto.SearchToolsResult{
				Score:        result.Score,
				ProjectID:    m.ProjectID,
				Name:         m.Name,
				Description:  m.Description,
				Technologies: m.Technologies,
				Domains:      m.Domains,
				Features:     m.Features,
			})
		} else {
			// Keep the original result if no metadata found
			results = append(results, result)
		}
	}

	// Sort results by score (highest first)
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	return results, nil
}

// GetSimilarProjects finds projects similar to the given project
func (s *HybridSearchService) GetSimilarProjects(ctx context.Context, projectID string) ([]dto.SearchToolsResult, error) {
	if s.vectorService == nil {
		return nil, fmt.Errorf("vector service not initialized")
	}
	return s.vectorService.GetSimilarProjects(ctx, projectID)
}

// SearchByTechnology searches for projects using a specific technology
func (s *HybridSearchService) SearchByTechnology(ctx context.Context, tech string) ([]dto.SearchToolsResult, error) {
	// Fall back to vector search with technology as query
	if s.graphService == nil {
		return s.VectorSearch(ctx, tech)
	}
	return s.graphService.SearchByTechnology(ctx, tech)
}

// SearchByDomain searches for projects in a specific domain
func (s *HybridSearchService) SearchByDomain(ctx context.Context, domain string) ([]dto.SearchToolsResult, error) {
	// Fall back to vector search with domain as query
	if s.graphService == nil {
		return s.VectorSearch(ctx, domain)
	}
	return s.graphService.SearchByDomain(ctx, domain)
}
