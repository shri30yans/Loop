package main

import (
	"Loop/auth"
	db "Loop/database"
	"Loop/projects"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	fmt.Println("Starting server...")
	db.StartServer()
	if secretKey := os.Getenv("JWT_SECRET"); secretKey != "" {
		auth.JwtSecret = []byte(secretKey)
	}
	router := mux.NewRouter()
	apiRouter := router.PathPrefix("/api").Subrouter()

	apiRouter.HandleFunc("/", db.Root).Methods("GET")

	// Auth routes /api/auth
	authRouter := apiRouter.PathPrefix("/auth").Subrouter()
	authRouter.HandleFunc("/register", auth.HandleRegister).Methods("POST")
	authRouter.HandleFunc("/login", auth.HandleLogin).Methods("POST")
	authRouter.HandleFunc("/verify", auth.HandleVerify).Methods("GET")

	// Project routes /api/project
	projectRouter := apiRouter.PathPrefix("/project").Subrouter()
	projectRouter.HandleFunc("/get_projects", auth.AuthMiddleware(projects.HandleGetProjects)).Methods("GET")
	projectRouter.HandleFunc("/get_project_info", auth.AuthMiddleware(projects.HandleGetProjectInfo)).Methods("GET")
	projectRouter.HandleFunc("/create_project", auth.AuthMiddleware(projects.HandleCreateProject)).Methods("POST")

	// Global middleware
	router.Use(corsMiddleware)

	log.Fatal(http.ListenAndServe(":8080", router))
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
