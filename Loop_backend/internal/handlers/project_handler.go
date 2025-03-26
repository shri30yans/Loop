package handlers

import (
	"Loop_backend/internal/dto"
	"Loop_backend/internal/middleware"
	"Loop_backend/internal/models"
	"Loop_backend/internal/response"
	"Loop_backend/internal/services"
	"net/http"

	"github.com/google/uuid"
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

func (h *ProjectHandler) RegisterRoutes(r RouteRegister) {
	r.RegisterProtectedRoute("/api/project/create", h.CreateProject, &dto.CreateProjectRequest{})
	r.RegisterProtectedRoute("/api/project/search", h.SearchProjects, nil)
	r.RegisterProtectedRoute("/api/project/{project_id:[a-fA-F0-9-]+}", h.GetProjectInfo, nil)
	r.RegisterProtectedRoute("/api/project/{project_id:[a-fA-F0-9-]+}/delete", h.DeleteProject, nil)
}

func (h *ProjectHandler) SearchProjects(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	keyword := vars["keyword"]

	var projects []*models.ProjectInfo
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

	// Validate UUID format
	if _, err := uuid.Parse(projectID); err != nil {
		response.RespondWithError(w, http.StatusBadRequest, "Invalid project ID format")
		return
	}

	project, err := h.projectService.GetProject(projectID)
	if err != nil {
		if err.Error() == "project not found" {
			response.RespondWithError(w, http.StatusNotFound, "Project not found")
			return
		}
		response.RespondWithError(w, http.StatusInternalServerError, "Error retrieving project")
		return
	}

	response.RespondWithJSON(w, http.StatusOK, project)
}

func (h *ProjectHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		response.RespondWithError(w, http.StatusBadRequest, "Invalid user ID format")
		return
	}

	req, ok := middleware.GetDTO[*dto.CreateProjectRequest](r)
	if !ok {
		response.RespondWithError(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	err = h.projectService.CreateProject(req, userUUID)
	if err != nil {
		response.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.RespondWithJSON(w, http.StatusCreated, map[string]string{
		"message": "Project created successfully",
	})
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
