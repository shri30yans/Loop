package handlers

import (
	"Loop_backend/internal/dto"
	"Loop_backend/internal/middleware"
	"Loop_backend/internal/response"
	"Loop_backend/internal/services"
	"fmt"
	"net/http"
	"strings"
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
	r.RegisterProtectedRoute("/api/project/search", h.HandleProjectSearch, nil) // Using query params
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

// HandleProjectSearch handles unified search requests from the frontend
func (h *SearchHandler) HandleProjectSearch(w http.ResponseWriter, r *http.Request) {
	// Extract keyword from query parameter
	keyword := r.URL.Query().Get("keyword")
	// Support mode parameter for future frontend toggle
	mode := r.URL.Query().Get("mode")

	var results []map[string]interface{}
	var err error

	if keyword == "" {
		// Fetch all projects when no keyword is provided
		fmt.Println("No keyword provided, fetching all projects")
		results, err = h.queryService.GetAllProjects()
		if err != nil {
			response.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch projects: "+err.Error())
			return
		}
	} else {
		// Determine search type based on mode parameter or auto-detection
		var useTopicSearch bool

		if mode == "advanced" {
			// Frontend explicitly requested advanced search
			useTopicSearch = true
		} else if mode == "basic" {
			// Frontend explicitly requested basic search
			useTopicSearch = false
		} else {
			// Auto-detect based on query pattern (hybrid approach)
			fmt.Println("Auto-detecting search type for query:", keyword)
			useTopicSearch = shouldUseTopicSearch(keyword)
		}

		if useTopicSearch {
			fmt.Printf("Using topic search for query: %s\n", keyword)
			results, err = h.queryService.ExecuteTopicSearchQuery(keyword)
		} else {
			fmt.Printf("Using regular search for query: %s\n", keyword)
			results, err = h.queryService.ExecuteSearchQuery(keyword)
		}

		if err != nil {
			response.RespondWithError(w, http.StatusInternalServerError, "Failed to execute search: "+err.Error())
			return
		}
	}

	// Format results for frontend
	projects := formatProjectResults(results)

	// Enrich projects with summaries
	enrichedProjects := make([]map[string]interface{}, 0, len(projects))
	for _, project := range projects {
		// Get the project ID
		projectID, ok := project["id"].(string)
		if !ok || projectID == "" {
			// If no valid ID, just add the project without summary
			enrichedProjects = append(enrichedProjects, project)
			continue
		}

		// Generate summary
		summary, err := h.summaryService.GenerateProjectSummary(projectID)
		if err != nil {
			fmt.Printf("Error generating summary for project %s: %v\n", projectID, err)
			// Add project without summary if generation fails
			enrichedProjects = append(enrichedProjects, project)
			continue
		}

		// If summary contains a 'summary' field, extract it, otherwise use whole summary
		var summaryText interface{}
		if summaryContent, ok := summary["summary"]; ok {
			summaryText = summaryContent
		} else {
			summaryText = summary
		}

		// Add summary to the project
		project["summary"] = summaryText
		enrichedProjects = append(enrichedProjects, project)
	}

	// Return in format expected by frontend
	response.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"projects": enrichedProjects,
		"total":    len(enrichedProjects),
	})
}

// Enhanced hybrid approach to determine which search type to use
func shouldUseTopicSearch(query string) bool {
	// Topic search indicators - phrases and words that suggest semantic search
	topicIndicators := []string{
		// Question formats
		"about", "related to", "concerning", "regarding", "similar to",
		"like", "such as", "involving", "dealing with", "associated with",

		// Command formats
		"find", "search", "show", "get", "retrieve", "locate",
		"list", "discover", "explore", "suggest", "recommend",

		// Topic phrases
		"projects on", "projects about", "projects using", "projects with",
		"technology for", "solution for", "approach to", "method for",
		"implementation of", "application of", "system for",

		// Domain-specific indicators
		"concept", "field", "domain", "area", "topic", "subject",
		"category", "theme", "focus", "discipline", "expertise",

		// Question words (if followed by topic)
		"what", "which", "where", "who", "how",
	}

	// Check for clear indicators first
	queryLower := strings.ToLower(query)
	for _, indicator := range topicIndicators {
		if strings.Contains(queryLower, indicator) {
			return true
		}
	}

	// Direct ID or slug search indicators (likely not topic search)
	directSearchIndicators := []string{
		"id:", "project:", "projectid:",
		"uuid:", "identifier:", "pid:",
	}

	for _, indicator := range directSearchIndicators {
		if strings.HasPrefix(queryLower, indicator) {
			return false
		}
	}

	// Short queries with single terms are likely direct searches
	if strings.Count(query, " ") <= 1 && len(query) < 20 {
		return false
	}

	// For queries that are 3+ words, use topic search
	words := strings.Fields(query)
	if len(words) >= 3 {
		return true
	}

	// Default to regular search for ambiguous cases
	return false
}

// Format search results to match the expected frontend format
func formatProjectResults(results []map[string]interface{}) []map[string]interface{} {
	formattedResults := make([]map[string]interface{}, 0, len(results))

	for _, result := range results {
		fmt.Printf("Processing result: %v\n", result)

		// IMPORTANT CHANGE: Extract project ID correctly
		var projectID string

		// Try to get projectId first
		if id, ok := result["projectId"].(string); ok && id != "" {
			projectID = id
		} else {
			// If that fails, try p.project_id
			if pID, ok := result["p.project_id"].(string); ok && pID != "" {
				projectID = pID
			} else {
				// Try other variations we've seen in the data
				for key, value := range result {
					if (strings.Contains(strings.ToLower(key), "project_id") ||
						strings.Contains(strings.ToLower(key), "projectid")) &&
						value != nil {
						if idStr, ok := value.(string); ok && idStr != "" {
							projectID = idStr
							break
						}
					}
				}
			}
		}

		// If we still couldn't find a valid ID, log and skip
		if projectID == "" {
			fmt.Println("Skipping project without ID:", result)
			continue
		}

		// Get project name - similarly try multiple field names
		var projectName string
		if name, ok := result["projectName"].(string); ok && name != "" {
			projectName = name
		} else if name, ok := result["p.name"].(string); ok && name != "" {
			projectName = name
		} else {
			projectName = "Unnamed Project"
		}

		// Get description - similarly flexible approach
		var description string
		if desc, ok := result["description"].(string); ok {
			description = desc
		} else if desc, ok := result["p.description"].(string); ok {
			description = desc
		}

		// Initialize tags as an empty SLICE, not nil
		tags := make([]string, 0)

		// Process tags if they exist
		if tagsInterface, ok := result["tags"]; ok && tagsInterface != nil {
			if tagsArray, ok := tagsInterface.([]string); ok {
				tags = tagsArray
			} else if tagsArray, ok := tagsInterface.([]interface{}); ok {
				for _, tag := range tagsArray {
					if tagStr, ok := tag.(string); ok {
						tags = append(tags, tagStr)
					}
				}
			}
		}

		// Fix for the status field
		status := "published" // Default value
		if statusVal, ok := result["status"].(string); ok {
			status = statusVal
		} else if statusVal, ok := result["p.status"].(string); ok {
			status = statusVal
		}

		// Format for frontend
		project := map[string]interface{}{
			"id":          projectID,
			"title":       projectName,
			"description": description,
			"tags":        tags,
			"status":      status,
		}

		formattedResults = append(formattedResults, project)
	}

	// Debug the final result before returning
	fmt.Printf("Formatted %d results\n", len(formattedResults))
	for i, p := range formattedResults {
		fmt.Printf("Project %d: ID=%v, Title=%v, Tags=%v\n", i, p["id"], p["title"], p["tags"])
	}

	return formattedResults
}
