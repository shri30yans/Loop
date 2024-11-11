package projects

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func HandleGetProjects(w http.ResponseWriter, r *http.Request) {
	projects, err := FetchProjects()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(projects)
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
	fmt.Println(newProject)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println()
	projectID, err := CreateProject(newProject.Title, newProject.Description, newProject.Introduction, newProject.Tags, newProject.OwnerID, newProject.Sections)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	newProject.ProjectID = projectID

	json.NewEncoder(w).Encode(newProject)
}
