package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"Loop_backend/config"
	"Loop_backend/internal/ai/providers"
	"Loop_backend/internal/handlers"
	authRepo "Loop_backend/internal/repositories/auth"
	graphRepo "Loop_backend/internal/repositories/graph"
	projectRepo "Loop_backend/internal/repositories/project"
	userRepo "Loop_backend/internal/repositories/user"
	"Loop_backend/internal/services/auth"
	"Loop_backend/internal/services/graph"
	rparser "Loop_backend/internal/services/graph/response_parser"
	"Loop_backend/internal/services/project"
	"Loop_backend/internal/services/user"
	neo4j "Loop_backend/platform/database/neo4j"
	postgres "Loop_backend/platform/database/postgres"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

type application struct {
	config         *config.Config
	authService    auth.AuthService
	userService    user.UserService
	projectService project.ProjectService
	authHandler    *handlers.AuthHandler
	userHandler    *handlers.UserHandler
	projectHandler *handlers.ProjectHandler
}

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found")
	}

	// Load Configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	app, err := initializeApp(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	// Initialize Router with gorilla/mux
	router := mux.NewRouter()
	routeRegister := handlers.NewRouteRegister(router, app.authService)

	// Register routes for all handlers
	app.authHandler.RegisterRoutes(routeRegister)
	app.userHandler.RegisterRoutes(routeRegister)
	app.projectHandler.RegisterRoutes(routeRegister)

	// Setup CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})

	// Wrap the router with the CORS handler
	handler := c.Handler(router)

	// Setup graceful shutdown
	defer postgres.Close()
	defer neo4j.Close()

	// Start Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	serverAddr := fmt.Sprintf(":%s", port)
	log.Printf("Server running at %s", serverAddr)

	server := &http.Server{
		Addr:    serverAddr,
		Handler: handler,
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed: %v", err)
	}
}

func initializeApp(cfg *config.Config) (*application, error) {
	if err := postgres.InitDB(&cfg.RelationalDatabaseConfig); err != nil {
		return nil, fmt.Errorf("failed to initialize database: %v", err)
	}
	if err := neo4j.InitNeo4j(&cfg.Neo4jConfig); err != nil {
		return nil, fmt.Errorf("failed to initialize graph database: %v", err)
	}

	// Get Database Instance
	db := postgres.GetDBPool()
	graphDB := neo4j.GetDriver()

	// Initialize AI Components
	provider, err := providers.NewProvider(&cfg.AIConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize AI provider: %v", err)
	}

	// Initialize Repositories
	uRepo := userRepo.NewUserRepository(db)
	pRepo := projectRepo.NewProjectRepository(db)
	aRepo := authRepo.NewAuthRepository(db)
	gRepo := graphRepo.NewGraphRepository(graphDB)

	// Initialize Services
	responseParser := rparser.NewResponseParser()
	graphService := graph.NewGraphService(gRepo, responseParser)

	projectProcessor := project.NewProjectProcessor(
		provider,
		graphService,
	)

	// Initialize Core Services
	authService := auth.NewAuthService(cfg.JWTConfig.Secret, aRepo)
	userService := user.NewUserService(uRepo)
	projectService := project.NewProjectService(pRepo, projectProcessor)

	// Initialize Handlers
	userHandler := handlers.NewUserHandler(userService)
	authHandler := handlers.NewAuthHandler(userService, authService)
	projectHandler := handlers.NewProjectHandler(projectService)

	return &application{
		config:         cfg,
		authService:    authService,
		userService:    userService,
		projectService: projectService,
		// tagService:     tagService,
		authHandler:    authHandler,
		userHandler:    userHandler,
		projectHandler: projectHandler,
	}, nil
}
