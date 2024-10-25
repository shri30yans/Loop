package handler

import (
	"Loop/auth"
	db "Loop/database"
	"Loop/projects"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

func StartServer() {
	err := db.InitDB()
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	defer db.DB.Close()
}

func Root(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Welcome to the Loop Backend API"))
}

func HandleGetProjects(w http.ResponseWriter, r *http.Request) {
	projects, err := projects.FetchProjects()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(projects)
}

func HandleCreateProject(w http.ResponseWriter, r *http.Request) {

	var newProject projects.Project
	err := json.NewDecoder(r.Body).Decode(&newProject)
	fmt.Println(newProject)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println()
	projectID, err := projects.CreateProject(newProject.Title, newProject.Description, newProject.Introduction, newProject.Tags)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	newProject.ProjectID = projectID
	json.NewEncoder(w).Encode(newProject)
}

func HandleRegister(w http.ResponseWriter, r *http.Request) {
	var req auth.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//if !auth.isValidEmail(req.Email) || !auth.isValidPassword(req.Password) {
	//	http.Error(w, "Invalid email or password", http.StatusBadRequest)
	//	return
	//}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	user, err := auth.CreateUser(req.Email, string(hashedPassword))
	if err != nil {
		//if auth.isDuplicateEmail(err) {
		//	http.Error(w, "Email already exists", http.StatusConflict)
		//	return
		//}
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	token, err := auth.GenerateToken(user.ID)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(auth.AuthResponse{
		Token: token,
		User:  user,
	})
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	var req auth.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := auth.GetUserByEmail(req.Email)
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := auth.GenerateToken(user.ID)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(auth.AuthResponse{
		Token: token,
		User:  user,
	})
}
