package projects

import (
	. "Loop/models"
	"encoding/json"
	"net/http"
	"strconv"
)

func HandleGetProjects(w http.ResponseWriter, r *http.Request) {
	keyword := r.URL.Query().Get("keyword")
	var projects []ProjectsResponse
	var total int
	var err error

	if keyword != "" {
		projects, total, err = FetchProjects(&keyword)
	} else {
		projects, total, err = FetchProjects(nil)
	}

	if err != nil {
		switch e := err.(type) {
		case *ErrNoProjects:
			http.Error(w, e.Error(), http.StatusNotFound)
		default:
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	// Create response with projects and total count
	response := map[string]interface{}{
		"total":    total,
		"projects": projects,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func HandleGetProjectInfo(w http.ResponseWriter, r *http.Request) {
	projectID := r.URL.Query().Get("project-id")
	if projectID == "" {
		http.Error(w, "Missing project-id parameter", http.StatusBadRequest)
		return
	}

	projectIDInt, err := strconv.Atoi(projectID)
	if err != nil {
		http.Error(w, "Invalid project-id parameter", http.StatusBadRequest)
		return
	}

	projects, err := FetchProjectInfo(projectIDInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(projects)
}

func HandleCreateProject(w http.ResponseWriter, r *http.Request) {

	var newProject Project
	err := json.NewDecoder(r.Body).Decode(&newProject)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//fmt.Println()
	projectID, err := CreateProject(newProject.Title, newProject.Description, newProject.Introduction, newProject.Tags, newProject.OwnerID, newProject.Sections)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	newProject.ProjectID = projectID

	json.NewEncoder(w).Encode(newProject)
}
