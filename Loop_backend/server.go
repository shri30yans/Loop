package main

import (
	"Loop_backend/database"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	fmt.Println("Starting server...")
	err := database.InitDB()
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	defer database.DB.Close()

	router := mux.NewRouter()
	router.HandleFunc("/", root).Methods("GET")
	router.HandleFunc("/projects", handleGetProjects).Methods("GET")
	router.HandleFunc("/create_project", handleCreateProject).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", router))
}

func root(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Welcome to the Loop Backend API"))
}

func handleGetProjects(w http.ResponseWriter, r *http.Request) {
	projects, err := database.FetchProjects()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(projects)
}

func handleCreateProject(w http.ResponseWriter, r *http.Request) {

	var newProject database.Project
	err := json.NewDecoder(r.Body).Decode(&newProject)
	fmt.Println(newProject)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println()
	projectID, err := database.CreateProject(newProject.Title, newProject.Description, newProject.Introduction, newProject.Tags)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	newProject.ProjectID = projectID
	json.NewEncoder(w).Encode(newProject)
}
