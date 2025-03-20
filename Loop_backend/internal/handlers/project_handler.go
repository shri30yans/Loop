package handlers

import (
	"Loop_backend/internal/dto"
	"Loop_backend/internal/middleware"
	"Loop_backend/internal/models"
	"Loop_backend/internal/response"
	"Loop_backend/internal/services"
	"encoding/json"
<<<<<<< HEAD
	"net/http"

	"github.com/gorilla/mux"
=======
	"fmt"
	"net/http"
>>>>>>> 4a2f436bed91636c5c2e3782993f5ab211ecfca7
)

type ProjectHandler struct {
	projectService services.ProjectService
}

func NewProjectHandler(projectService services.ProjectService) *ProjectHandler {
	return &ProjectHandler{
		projectService: projectService,
	}
}

<<<<<<< HEAD
func (h *ProjectHandler) RegisterRoutes(r *RouteRegister) {
	r.RegisterProtectedRoute("/api/project/search", h.SearchProjects)
	r.RegisterProtectedRoute("/api/project/{project_id:[a-fA-F0-9-]+}", h.GetProjectInfo)
	r.RegisterProtectedRoute("/api/project/create", h.CreateProject)
	r.RegisterProtectedRoute("/api/project/{uuid}/delete", h.DeleteProject)
=======
func (h *ProjectHandler) GetAllProjects(w http.ResponseWriter, r *http.Request) {
	projects, err := h.projectService.GetProjects()
	if err != nil {
		fmt.Println("here with error")
		response.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"projects": projects,
		"total":    len(projects),
	})
>>>>>>> 4a2f436bed91636c5c2e3782993f5ab211ecfca7
}

func (h *ProjectHandler) SearchProjects(w http.ResponseWriter, r *http.Request) {
	keyword := r.URL.Query().Get("keyword")

	var projects []*models.Project
<<<<<<< HEAD
	var err error

	projects, err = h.projectService.SearchProjects(keyword)
=======
	var count int
	var err error

	projects, count, err = h.projectService.SearchProjects(keyword)
>>>>>>> 4a2f436bed91636c5c2e3782993f5ab211ecfca7

	if err != nil {
		response.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"projects": projects,
<<<<<<< HEAD
=======
		"total":    count,
>>>>>>> 4a2f436bed91636c5c2e3782993f5ab211ecfca7
	})
}

func (h *ProjectHandler) GetProjectInfo(w http.ResponseWriter, r *http.Request) {
<<<<<<< HEAD
	vars := mux.Vars(r)
	projectID := vars["project_id"]
=======
	projectID := r.URL.Query().Get("project-id")
>>>>>>> 4a2f436bed91636c5c2e3782993f5ab211ecfca7

	project, err := h.projectService.GetProject(projectID)
	if err != nil {
		response.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.RespondWithJSON(w, http.StatusOK, project)
}

func (h *ProjectHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)

<<<<<<< HEAD
	// Convert to dto.CreateProjectRequest
	var req dto.CreateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
=======
	var req dto.CreateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
>>>>>>> 4a2f436bed91636c5c2e3782993f5ab211ecfca7
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
<<<<<<< HEAD
	vars := mux.Vars(r)
	projectID := vars["project_id"]
=======
	projectID := r.URL.Query().Get("project-id")
>>>>>>> 4a2f436bed91636c5c2e3782993f5ab211ecfca7

	if err := h.projectService.DeleteProject(projectID); err != nil {
		response.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.RespondWithJSON(w, http.StatusOK, map[string]string{
		"message": "Project deleted successfully",
	})
}
<<<<<<< HEAD
=======

func (h *ProjectHandler) RegisterRoutes(r *RouteRegister) {
	r.RegisterProtectedRoute("/api/project/get_projects", h.GetAllProjects) // Get all projects
	r.RegisterProtectedRoute("/api/project/search", h.SearchProjects)       // Search projects
	r.RegisterProtectedRoute("/api/project/info", h.GetProjectInfo)
	r.RegisterProtectedRoute("/api/project/create", h.CreateProject)
	r.RegisterProtectedRoute("/api/project/delete", h.DeleteProject)
}
>>>>>>> 4a2f436bed91636c5c2e3782993f5ab211ecfca7
