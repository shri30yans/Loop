// main.go
package main

import (
	"Loop/auth"
	db "Loop/database"
	"Loop/projects"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	fmt.Println("Starting server...")
	db.StartServer()
	router := mux.NewRouter()
	apiRouter := router.PathPrefix("/api").Subrouter()

	// Apply CORS middleware to all routes
	router.Use(corsMiddleware)

	apiRouter.HandleFunc("/", db.Root).Methods("GET", "OPTIONS")

	userRouter := apiRouter.PathPrefix("/user").Subrouter()
	userRouter.HandleFunc("/get_user_info", auth.HandleGetUserInfo).Methods("GET", "OPTIONS")

	// Auth routes /api/auth
	authRouter := apiRouter.PathPrefix("/auth").Subrouter()
	authRouter.HandleFunc("/register", auth.HandleRegister).Methods("POST", "OPTIONS")
	authRouter.HandleFunc("/login", auth.HandleLogin).Methods("POST", "OPTIONS")
	authRouter.HandleFunc("/change_password", auth.HandleChangePassword).Methods("POST", "OPTIONS")
	authRouter.HandleFunc("/verify", auth.HandleVerify).Methods("GET", "OPTIONS")

	// Project routes /api/project
	projectRouter := apiRouter.PathPrefix("/project").Subrouter()
	projectRouter.HandleFunc("/get_projects", auth.AuthMiddleware(projects.HandleGetProjects)).Methods("GET", "OPTIONS")
	projectRouter.HandleFunc("/get_project_info", auth.AuthMiddleware(projects.HandleGetProjectInfo)).Methods("GET", "OPTIONS")
	projectRouter.HandleFunc("/create_project", auth.AuthMiddleware(projects.HandleCreateProject)).Methods("POST", "OPTIONS")

	log.Fatal(http.ListenAndServe(":8080", router))
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "86400") // 24 hours

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
