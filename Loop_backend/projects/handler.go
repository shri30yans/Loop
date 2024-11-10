package projects

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func HandleGetProjects(w http.ResponseWriter, r *http.Request) {
	projects, err := FetchProjects()
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
