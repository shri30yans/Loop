package handlers

import (
	"Loop_backend/internal/dto"
	"Loop_backend/internal/middleware"
	"Loop_backend/internal/response"
	"Loop_backend/internal/services/search"
	"context"
	"net/http"
)

// SearchHandler handles search-related routes
type SearchHandler struct {
	searchService search.SearchService
}

// NewSearchHandler creates a new search handler
func NewSearchHandler(searchService search.SearchService) *SearchHandler {
	return &SearchHandler{
		searchService: searchService,
	}
}

// RegisterRoutes registers all search routes
func (h *SearchHandler) RegisterRoutes(router RouteRegister) {
	// Single endpoint for all search types
	router.RegisterPublicRouteWithMethod("/api/search_tools", http.MethodPost, middleware.ValidateRequest(h.HandleSearchTools, &dto.SearchToolsRequest{}))
}

// SearchToolsRequest represents the unified search request
type SearchToolsRequest struct {
	Query string `json:"query" validate:"required"`
	Mode  string `json:"mode"` // "hybrid", "vector", or "graph"
	Limit int    `json:"limit,omitempty"`
}

// HandleSearchTools handles all search requests through a single endpoint
func (h *SearchHandler) HandleSearchTools(w http.ResponseWriter, r *http.Request) {
	req, ok := middleware.GetDTO[*dto.SearchToolsRequest](r)
	if !ok {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Query == "" {
		response.Error(w, http.StatusBadRequest, "Query is required")
		return
	}

	ctx := context.Background()
	var results []dto.SearchToolsResult
	var err error

	// TODO: Ensure searchService has functional HybridSearch and GraphSearch methods
	switch req.Mode { 
	case "vector":
		results, err = h.searchService.VectorSearch(ctx, req.Query)
	case "graph":
		results, err = h.searchService.GraphSearch(ctx, req.Query)
	case "hybrid": // Added explicit hybrid case
		results, err = h.searchService.HybridSearch(ctx, req.Query)
	default:
		// Default to hybrid search if mode is not specified or is unknown.
		// Consider logging a warning for an unrecognized mode.
		// log.Printf("Warning: Unrecognized search mode '%s', defaulting to hybrid search", req.Mode)
		results, err = h.searchService.HybridSearch(ctx, req.Query)
	}
	// results, err = h.searchService.VectorSearch(ctx, req.Query) // Original hardcoded call removed

	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to perform search")
		return
	}

	// Apply limit if specified
	if req.Limit > 0 && len(results) > req.Limit {
		results = results[:req.Limit]
	}

	resp := dto.SearchResponse{
		Results: results,
		Total:   len(results),
	}

	response.JSON(w, http.StatusOK, resp)
}
