package handlers

import (
	"Loop_backend/internal/dto"
	"Loop_backend/internal/middleware"
	"Loop_backend/internal/response"
	"Loop_backend/internal/services"
	"fmt"
	"net/http"
)

// SearchHandler handles search-related requests
type SearchHandler struct {
	queryService   services.QueryService
	summaryService services.SummaryService
}

// NewSearchHandler creates a new search handler instance
func NewSearchHandler(
	queryService services.QueryService,
	summaryService services.SummaryService,
) *SearchHandler {
	return &SearchHandler{
		queryService:   queryService,
		summaryService: summaryService,
	}
}

// RegisterRoutes registers all routes for the search handler
func (h *SearchHandler) RegisterRoutes(r RouteRegister) {
	r.RegisterProtectedRoute("/api/search", h.HandleSearch, &dto.SearchRequest{})
	r.RegisterProtectedRoute("/api/topic-search", h.HandleTopicSearch, &dto.SearchRequest{})
}

// HandleSearch converts natural language to cypher query and returns results
func (h *SearchHandler) HandleSearch(w http.ResponseWriter, r *http.Request) {
	// Get the request from context (parsed by middleware)
	req, ok := middleware.GetDTO[*dto.SearchRequest](r)
	if !ok {
		response.RespondWithErrorDetails(w, http.StatusBadRequest, "Invalid request payload", map[string]string{
			"reason":          "Failed to parse or validate request body",
			"expected_fields": "query",
		})
		return
	}

	// Transform the query to Cypher
	cypherQuery, err := h.queryService.TransformQueryToCypher(req.Query)
	if err != nil {
		response.RespondWithError(w, http.StatusInternalServerError, "Failed to transform query: "+err.Error())
		return
	}

	// Execute the query
	results, err := h.queryService.ExecuteSearchQuery(req.Query)
	if err != nil {
		response.RespondWithError(w, http.StatusInternalServerError, "Failed to execute query: "+err.Error())
		return
	}

	// Return the results
	response.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"cypher":  cypherQuery,
		"results": results,
	})
}

// HandleTopicSearch finds projects by topic and includes summaries
func (h *SearchHandler) HandleTopicSearch(w http.ResponseWriter, r *http.Request) {
	// Parse request
	req, ok := middleware.GetDTO[*dto.SearchRequest](r)
	if !ok {
		response.RespondWithErrorDetails(w, http.StatusBadRequest, "Invalid request payload", map[string]string{
			"reason": "Failed to parse request body",
		})
		return
	}

	// Transform query to get Cypher (for debugging/display)
	cypherQuery, err := h.queryService.TransformQueryToTopicCypher(req.Query)
	if err != nil {
		response.RespondWithError(w, http.StatusInternalServerError, "Failed to transform query: "+err.Error())
		return
	}

	// Execute with the specialized function instead of regular search
	results, err := h.queryService.ExecuteTopicSearchQuery(req.Query)
	if err != nil {
		response.RespondWithError(w, http.StatusInternalServerError, "Failed to execute query: "+err.Error())
		return
	}

	// Process results and get summaries
	enrichedResults := make([]map[string]interface{}, 0)

	for _, result := range results {
		projectID, ok := result["projectId"].(string)
		if !ok || projectID == "" {
			continue
		}

		summary, err := h.summaryService.GenerateProjectSummary(projectID)
		if err != nil {
			fmt.Printf("Error generating summary for project %s: %v\n", projectID, err)
			continue
		}

		enrichedResults = append(enrichedResults, summary)
	}

	response.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"query":    req.Query,
		"cypher":   cypherQuery,
		"projects": enrichedResults,
	})
}
