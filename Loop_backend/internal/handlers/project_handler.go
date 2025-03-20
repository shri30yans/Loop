package handlers

import (
	"Loop_backend/internal/dto"
	"Loop_backend/internal/middleware"
	"Loop_backend/internal/models"
	"Loop_backend/internal/response"
	"Loop_backend/internal/services"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type ProjectHandler struct {
	projectService services.ProjectService
}

func NewProjectHandler(projectService services.ProjectService) *ProjectHandler {
	return &ProjectHandler{
		projectService: projectService,
	}
}

func (h *ProjectHandler) RegisterRoutes(r *RouteRegister) {
	r.RegisterProtectedRoute("/api/project/search", h.SearchProjects)
	r.RegisterProtectedRoute("/api/project/{project_id:[a-fA-F0-9-]+}", h.GetProjectInfo)
	r.RegisterProtectedRoute("/api/project/create", h.CreateProject)
	r.RegisterProtectedRoute("/api/project/{uuid}/delete", h.DeleteProject)
}

func (h *ProjectHandler) SearchProjects(w http.ResponseWriter, r *http.Request) {
	keyword := r.URL.Query().Get("keyword")

	var projects []*models.Project
	var err error

	projects, err = h.projectService.SearchProjects(keyword)

	if err != nil {
		response.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"projects": projects,
	})
}

func (h *ProjectHandler) GetProjectInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectID := vars["project_id"]

	project, err := h.projectService.GetProject(projectID)
	if err != nil {
		response.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.RespondWithJSON(w, http.StatusOK, project)
}

func (h *ProjectHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)

	// Convert to dto.CreateProjectRequest
	var req dto.CreateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
	}

	req.OwnerID = userID

	project, err := h.projectService.CreateProject(req)
	if err != nil {
		response.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.RespondWithJSON(w, http.StatusCreated, project)
}

func (h *ProjectHandler) DeleteProject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectID := vars["project_id"]

	if err := h.projectService.DeleteProject(projectID); err != nil {
		response.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.RespondWithJSON(w, http.StatusOK, map[string]string{
		"message": "Project deleted successfully",
	})
}
