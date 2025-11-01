package handlers

import (
	"Loop_backend/internal/dto"
	"Loop_backend/internal/middleware"
	"Loop_backend/internal/response"
	"Loop_backend/internal/services/search"
	"context"
	"net/http"
)

// ToolsHandler handles tool-related routes
type ToolsHandler struct {
	searchService search.SearchService
	apiTokenMiddleware *middleware.APITokenMiddleware
}

// NewToolsHandler creates a new tools handler
func NewToolsHandler(searchService search.SearchService, apiToken string) *ToolsHandler {
	return &ToolsHandler{
		searchService: searchService,
		apiTokenMiddleware: middleware.NewAPITokenMiddleware(apiToken),
	}
}

// RegisterRoutes registers all tool routes
func (h *ToolsHandler) RegisterRoutes(router RouteRegister) {
	// Register tool endpoints with API token validation
	
	// Search tools endpoints - more granular API
	router.RegisterPublicRouteWithMethod(
		"/api/tools/search/vector", 
		http.MethodPost, 
		h.apiTokenMiddleware.ValidateAPIToken(
			middleware.ValidateRequest(h.HandleVectorSearch, &dto.SearchToolsRequest{}),
		),
	)
	
	router.RegisterPublicRouteWithMethod(
		"/api/tools/search/graph", 
		http.MethodPost, 
		h.apiTokenMiddleware.ValidateAPIToken(
			middleware.ValidateRequest(h.HandleGraphSearch, &dto.SearchToolsRequest{}),
		),
	)
	
	router.RegisterPublicRouteWithMethod(
		"/api/tools/search/hybrid", 
		http.MethodPost, 
		h.apiTokenMiddleware.ValidateAPIToken(
			middleware.ValidateRequest(h.HandleHybridSearch, &dto.SearchToolsRequest{}),
		),
	)
	
	// Backward compatibility - maintain the old endpoint but with token validation
	router.RegisterPublicRouteWithMethod(
		"/api/search_tools", 
		http.MethodPost, 
		h.apiTokenMiddleware.ValidateAPIToken(
			middleware.ValidateRequest(h.HandleSearchTools, &dto.SearchToolsRequest{}),
		),
	)
}

// HandleVectorSearch handles vector search requests
func (h *ToolsHandler) HandleVectorSearch(w http.ResponseWriter, r *http.Request) {
	req, ok := middleware.GetDTO[*dto.SearchToolsRequest](r)
	if !ok {
		appErr := response.NewValidationError("Invalid request body")
		response.RespondWithAppError(w, appErr)
		return
	}

	if req.Query == "" {
		appErr := response.NewValidationError("Query is required")
		response.RespondWithAppError(w, appErr)
		return
	}

	ctx := context.Background()
	results, err := h.searchService.VectorSearch(ctx, req.Query)
	if err != nil {
		appErr := response.NewSearchError("Failed to perform vector search")
		appErr.WithDetails(map[string]interface{}{
			"query": req.Query,
			"error": err.Error(),
		})
		response.RespondWithAppError(w, appErr)
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

// HandleGraphSearch handles graph search requests
func (h *ToolsHandler) HandleGraphSearch(w http.ResponseWriter, r *http.Request) {
	req, ok := middleware.GetDTO[*dto.SearchToolsRequest](r)
	if !ok {
		appErr := response.NewValidationError("Invalid request body")
		response.RespondWithAppError(w, appErr)
		return
	}

	if req.Query == "" {
		appErr := response.NewValidationError("Query is required")
		response.RespondWithAppError(w, appErr)
		return
	}

	ctx := context.Background()
	results, err := h.searchService.GraphSearch(ctx, req.Query)
	if err != nil {
		appErr := response.NewSearchError("Failed to perform graph search")
		appErr.WithDetails(map[string]interface{}{
			"query": req.Query,
			"error": err.Error(),
		})
		response.RespondWithAppError(w, appErr)
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

// HandleHybridSearch handles hybrid search requests
func (h *ToolsHandler) HandleHybridSearch(w http.ResponseWriter, r *http.Request) {
	req, ok := middleware.GetDTO[*dto.SearchToolsRequest](r)
	if !ok {
		appErr := response.NewValidationError("Invalid request body")
		response.RespondWithAppError(w, appErr)
		return
	}

	if req.Query == "" {
		appErr := response.NewValidationError("Query is required")
		response.RespondWithAppError(w, appErr)
		return
	}

	ctx := context.Background()
	results, err := h.searchService.HybridSearch(ctx, req.Query)
	if err != nil {
		appErr := response.NewSearchError("Failed to perform hybrid search")
		appErr.WithDetails(map[string]interface{}{
			"query": req.Query,
			"error": err.Error(),
		})
		response.RespondWithAppError(w, appErr)
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

// HandleSearchTools handles all search requests through a single endpoint (for backward compatibility)
func (h *ToolsHandler) HandleSearchTools(w http.ResponseWriter, r *http.Request) {
	req, ok := middleware.GetDTO[*dto.SearchToolsRequest](r)
	if !ok {
		appErr := response.NewValidationError("Invalid request body")
		response.RespondWithAppError(w, appErr)
		return
	}

	if req.Query == "" {
		appErr := response.NewValidationError("Query is required")
		response.RespondWithAppError(w, appErr)
		return
	}

	ctx := context.Background()
	var results []dto.SearchToolsResult
	var err error

	switch req.Mode {
	case "vector":
		results, err = h.searchService.VectorSearch(ctx, req.Query)
		if err != nil {
			appErr := response.NewSearchError("Failed to perform vector search")
			appErr.WithDetails(map[string]interface{}{
				"query": req.Query,
				"mode":  "vector",
				"error": err.Error(),
			})
			response.RespondWithAppError(w, appErr)
			return
		}
	case "graph":
		results, err = h.searchService.GraphSearch(ctx, req.Query)
		if err != nil {
			appErr := response.NewSearchError("Failed to perform graph search")
			appErr.WithDetails(map[string]interface{}{
				"query": req.Query,
				"mode":  "graph",
				"error": err.Error(),
			})
			response.RespondWithAppError(w, appErr)
			return
		}
	case "hybrid", "":
		// Default to hybrid search if mode is not specified or is "hybrid"
		results, err = h.searchService.HybridSearch(ctx, req.Query)
		if err != nil {
			appErr := response.NewSearchError("Failed to perform hybrid search")
			appErr.WithDetails(map[string]interface{}{
				"query": req.Query,
				"mode":  "hybrid",
				"error": err.Error(),
			})
			response.RespondWithAppError(w, appErr)
			return
		}
	default:
		appErr := response.NewValidationError("Invalid search mode")
		appErr.WithDetails(map[string]interface{}{
			"supported_modes": []string{"vector", "graph", "hybrid"},
			"provided_mode":   req.Mode,
		})
		response.RespondWithAppError(w, appErr)
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