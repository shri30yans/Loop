package main

import (
	auth "Loop/auth"
	"Loop/handler"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	fmt.Println("Starting server...")
	handler.StartServer()
	if secretKey := os.Getenv("JWT_SECRET"); secretKey != "" {
		auth.JwtSecret = []byte(secretKey)
	}
	router := mux.NewRouter()
	apiRouter := router.PathPrefix("/api").Subrouter()

	apiRouter.HandleFunc("/", handler.Root).Methods("GET")

	// Auth routes /api/auth
	authRouter := apiRouter.PathPrefix("/auth").Subrouter()
	authRouter.HandleFunc("/register", handler.HandleRegister).Methods("POST")
	authRouter.HandleFunc("/login", handler.HandleLogin).Methods("POST")

	// Project routes /api/project
	projectRouter := apiRouter.PathPrefix("/project").Subrouter()
	projectRouter.HandleFunc("/fetch_projects", auth.AuthMiddleware(handler.HandleGetProjects)).Methods("GET")
	projectRouter.HandleFunc("/create_project", auth.AuthMiddleware(handler.HandleCreateProject)).Methods("POST")

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
