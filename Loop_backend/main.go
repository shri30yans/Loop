package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"Loop_backend/config"
	"Loop_backend/internal/ai/ollama"
	"Loop_backend/internal/ai/processor"
	"Loop_backend/internal/handlers"
	"Loop_backend/internal/middleware"
	"Loop_backend/internal/repositories"
	"Loop_backend/internal/services"
	"Loop_backend/internal/services/tags"
	"Loop_backend/platform/database"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found")
	}

	// Initialize configuration
	cfg := config.New()

	// Set up PostgreSQL connection
	postgresDB, err := database.NewPostgresConnection(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer postgresDB.Close()

	// Set up Neo4j connection
	neo4jDriver, err := database.NewNeo4jConnection(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to Neo4j: %v", err)
	}
	defer func() {
		if err := neo4jDriver.Close(); err != nil {
			log.Printf("Error closing Neo4j connection: %v", err)
		}
	}()

	// Initialize router
	router := mux.NewRouter()

	// Initialize repositories
	projectRepo := repositories.NewPostgresProjectRepository(postgresDB)
	graphRepo := repositories.NewGraphRepository(neo4jDriver)
	entityRepo := repositories.NewPostgresEntityRepository(postgresDB)

	// Initialize AI services
	ollamaProvider := ollama.NewProvider("http://localhost:11434", "llama2")
	entityProcessor := processor.NewEntityProcessor(ollamaProvider)

	// Initialize services
	tagService := tags.NewTagGenerationService(projectRepo, ollamaProvider)
	entityProcessingSvc := services.NewEntityProcessingService(entityProcessor, graphRepo, entityRepo)
	projectService := services.NewProjectService(projectRepo, graphRepo, tagService, entityProcessingSvc)

	// Initialize handlers
	projectHandler := handlers.NewProjectHandler(projectService)

	// Health check route
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	// API routes
	api := router.PathPrefix("/api").Subrouter()
	api.Use(middleware.JSONMiddleware)

	// Project routes
	projects := api.PathPrefix("/projects").Subrouter()
	projects.HandleFunc("", projectHandler.CreateProject).Methods("POST")
	projects.HandleFunc("/{project_id}", projectHandler.GetProject).Methods("GET")
	projects.HandleFunc("/{project_id}", projectHandler.DeleteProject).Methods("DELETE")
	projects.HandleFunc("/search", projectHandler.SearchProjects).Methods("GET")

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), router); err != nil {
		log.Fatal(err)
	}
}
