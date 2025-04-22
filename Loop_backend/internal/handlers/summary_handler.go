package handlers

import (
	"Loop_backend/internal/response"
	"Loop_backend/internal/services"
	"net/http"

	"github.com/gorilla/mux"
)

// SummaryHandler handles project summary requests
type SummaryHandler struct {
	summaryService services.SummaryService
}

// NewSummaryHandler creates a new summary handler
func NewSummaryHandler(summaryService services.SummaryService) *SummaryHandler {
	return &SummaryHandler{
		summaryService: summaryService,
	}
}

// RegisterRoutes registers all routes for the summary handler
func (h *SummaryHandler) RegisterRoutes(r RouteRegister) {
	r.RegisterProtectedRoute("/api/projects/{id}/summary", h.GetProjectSummary, nil)
}

// GetProjectSummary handles the request to generate a project summary
func (h *SummaryHandler) GetProjectSummary(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectID := vars["id"]

	if projectID == "" {
		response.RespondWithError(w, http.StatusBadRequest, "Missing project ID")
		return
	}

	summary, err := h.summaryService.GenerateProjectSummary(projectID)
	if err != nil {
		response.RespondWithError(w, http.StatusInternalServerError, "Failed to generate summary: "+err.Error())
		return
	}

	response.RespondWithJSON(w, http.StatusOK, summary)
}
