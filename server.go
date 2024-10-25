package main

import (
	"Loop_backend/database"
	"Loop_backend/auth"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var authenticator *auth.Authenticator

func main() {
	fmt.Println("Starting server...")
	
	// Initialize Auth0
	var err error
	authenticator, err = auth.NewAuthenticator()
	if err != nil {
		log.Fatalf("Error initializing authenticator: %v", err)
	}

	// Initialize database
	err = database.InitDB()
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	defer database.DB.Close()

	router := mux.NewRouter()
	
	// Public routes
	router.HandleFunc("/", root).Methods("GET")
	router.HandleFunc("/login", handleLogin).Methods("GET")
	router.HandleFunc("/callback", handleCallback).Methods("GET")

	// Protected routes
	router.HandleFunc("/projects", auth.AuthMiddleware(handleGetProjects)).Methods("GET")
	router.HandleFunc("/create_project", auth.AuthMiddleware(handleCreateProject)).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", router))
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	state := generateRandomState()
	// Store state in session or cookie
	url := authenticator.Config.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func handleCallback(w http.ResponseWriter, r *http.Request) {
	// Handle Auth0 callback
	code := r.URL.Query().Get("code")
	token, err := authenticator.Config.Exchange(r.Context(), code)
	if err != nil {
		http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
		return
	}

	// Get user info from Auth0
	userInfo, err := getUserInfo(token.AccessToken)
	if err != nil {
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}

	// Create or update user in your database
	// Return JWT token or set session
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userInfo)
}

// ... rest of your existing handlers ...
