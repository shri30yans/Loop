package search

import (
	"Loop_backend/internal/dto"
	"context"
)

// VectorService handles vector-based search operations
type VectorService interface {
	SearchProjects(ctx context.Context, query string) ([]dto.SearchToolsResult, error)
	GetSimilarProjects(ctx context.Context, projectID string) ([]dto.SearchToolsResult, error)
}

// GraphService handles graph-based search operations
type GraphService interface {
	// Basic graph search operations
	SearchProjects(ctx context.Context, query string) ([]dto.SearchToolsResult, error)

	// Specialized graph queries
	SearchByTechnology(ctx context.Context, tech string) ([]dto.SearchToolsResult, error)
	SearchByDomain(ctx context.Context, domain string) ([]dto.SearchToolsResult, error)

	// Query transformation
	TransformQueryToCypher(query string) (string, error)
}

// SearchService defines the unified interface for search operations
type SearchService interface {
	// Main search methods
	VectorSearch(ctx context.Context, query string) ([]dto.SearchToolsResult, error)
	GraphSearch(ctx context.Context, query string) ([]dto.SearchToolsResult, error)
	HybridSearch(ctx context.Context, query string) ([]dto.SearchToolsResult, error)

	// Additional search operations
	GetSimilarProjects(ctx context.Context, projectID string) ([]dto.SearchToolsResult, error)
	SearchByTechnology(ctx context.Context, tech string) ([]dto.SearchToolsResult, error)
	SearchByDomain(ctx context.Context, domain string) ([]dto.SearchToolsResult, error)
}

// Assert that HybridSearchService implements SearchService
var _ SearchService = (*HybridSearchService)(nil)
