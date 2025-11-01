package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
    "fmt"

	"Loop_backend/internal/dto"
	"Loop_backend/internal/middleware"
	"Loop_backend/internal/response"
	"Loop_backend/internal/services/project"
)

// ProjectHandler handles project-related requests
type ProjectHandler struct {
    projectService project.ProjectService
}

// NewProjectHandler creates a new project handler
func NewProjectHandler(projectService project.ProjectService, authMiddleware func(http.HandlerFunc) http.HandlerFunc) *ProjectHandler {
    return &ProjectHandler{
        projectService: projectService,
    }
}

// RegisterRoutes registers all routes for the project handler
func (h *ProjectHandler) RegisterRoutes(r RouteRegister) {
    // Register routes with specific HTTP methods
    r.RegisterProtectedRouteWithMethod("/api/project/search", http.MethodGet, h.ListProjects)
    r.RegisterProtectedRouteWithMethod("/api/project/create", http.MethodPost, middleware.ValidateRequest(h.CreateProject, &dto.CreateProjectRequest{}))
    r.RegisterProtectedRouteWithMethod("/api/project", http.MethodGet, h.ListProjects)
    r.RegisterProtectedRouteWithMethod("/api/project/{id}", http.MethodDelete, h.DeleteProject)
    r.RegisterProtectedRouteWithMethod("/api/project/{id}", http.MethodGet, h.GetProject)
}

// ListProjects handles the request to list all projects
func (h *ProjectHandler) ListProjects(w http.ResponseWriter, r *http.Request) {
    projects, err := h.projectService.ListProjects()
    if err != nil {
        response.RespondWithError(w, http.StatusInternalServerError, "Failed to list projects: "+err.Error())
        return
    }

    response.RespondWithJSON(w, http.StatusOK, projects)
}

// getUserIDFromContext extracts and validates the user ID from the context
func (h *ProjectHandler) getUserIDFromContext(w http.ResponseWriter, r *http.Request) (uuid.UUID, bool) {
    userIDStr, ok := r.Context().Value(middleware.UserIDKey).(string)
    if !ok {
        response.RespondWithError(w, http.StatusUnauthorized, "User ID not found in context")
        return uuid.Nil, false
    }

    userID, err := uuid.Parse(userIDStr)
    if err != nil {
        response.RespondWithError(w, http.StatusInternalServerError, "Failed to parse user ID: "+err.Error())
        return uuid.Nil, false
    }

    return userID, true
}

// GetProject handles the request to get a specific project
func (h *ProjectHandler) GetProject(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    projectID := vars["id"]

    // Parse UUID
    id, err := uuid.Parse(projectID)
    if err != nil {
        response.RespondWithError(w, http.StatusBadRequest, "Invalid project ID")
        return
    }

     // Get user ID from context
    // userID, ok := h.getUserIDFromContext(w, r)
    // if !ok {
    //     return
    // }

    project, err := h.projectService.GetProject(id)
    if err != nil {
        response.RespondWithError(w, http.StatusInternalServerError, "Failed to get project: "+err.Error())
        return
    }

    response.RespondWithJSON(w, http.StatusOK, project)
}

// CreateProject handles the request to create a new project
func (h *ProjectHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
    req, ok := middleware.GetDTO[*dto.CreateProjectRequest](r)
    if !ok {
        response.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
        return
    }

    // Get user ID from context
    userID, ok := h.getUserIDFromContext(w, r)
    if !ok {
        return
    }
    fmt.Printf("Starting project creation by user ID: %s\n", userID)

    // Convert DTO to model with owner ID
    project := req.ToModel(userID)

    if err := h.projectService.CreateProject(project); err != nil {
        response.RespondWithError(w, http.StatusInternalServerError, "Failed to create project: "+err.Error())
        return
    }

    response.RespondWithJSON(w, http.StatusCreated, project)
}

// DeleteProject handles the request to delete a project
func (h *ProjectHandler) DeleteProject(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    projectID := vars["id"]

    // Parse project UUID
    id, err := uuid.Parse(projectID)
    if err != nil {
        response.RespondWithError(w, http.StatusBadRequest, "Invalid project ID")
        return
    }

    // Get user ID from context
    userID, ok := h.getUserIDFromContext(w, r)
    if !ok {
        return
    }

    // Get project to verify ownership
    project, err := h.projectService.GetProject(id)
    if err != nil {
        response.RespondWithError(w, http.StatusInternalServerError, "Failed to get project: "+err.Error())
        return
    }

    // Verify project ownership
    if project.OwnerID != userID {
        response.RespondWithError(w, http.StatusForbidden, "You do not have permission to delete this project")
        return
    }

    if err := h.projectService.DeleteProject(id); err != nil {
        response.RespondWithError(w, http.StatusInternalServerError, "Failed to delete project: "+err.Error())
        return
    }

    response.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Project deleted successfully"})
}
