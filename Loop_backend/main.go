package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"Loop_backend/config"
	"Loop_backend/internal/ai/providers"
	"Loop_backend/internal/handlers"
	"Loop_backend/internal/repositories"
	"Loop_backend/internal/services"
	neo4j "Loop_backend/platform/database/neo4j"
	postgres "Loop_backend/platform/database/postgres"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

type application struct {
	config         *config.Config
	authService    services.AuthService
	userService    services.UserService
	projectService services.ProjectService
	queryService   services.QueryService
	authHandler    *handlers.AuthHandler
	userHandler    *handlers.UserHandler
	projectHandler *handlers.ProjectHandler
	searchHandler  *handlers.SearchHandler
	summaryService services.SummaryService
	summaryHandler *handlers.SummaryHandler
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
	app.searchHandler.RegisterRoutes(routeRegister)
	app.summaryHandler.RegisterRoutes(routeRegister)

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
	userRepo := repositories.NewUserRepository(db)
	projectRepo := repositories.NewProjectRepository(db)
	authRepo := repositories.NewAuthRepository(db)
	graphRepo := repositories.NewGraphRepository(graphDB)
	tagRepo := repositories.NewTagRepository(db)

	// Initialize Services
	graphService, err := services.NewGraphService(graphRepo)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize graph service: %w", err)
	}

	tagService := services.NewTagService(tagRepo, graphRepo)
	projectProcessor := services.NewProjectProcessor(
		provider,
		graphService,
		tagService,
	)

	queryService := services.NewQueryService(provider, graphRepo)
	summaryService := services.NewSummaryService(provider, projectRepo, graphRepo)

	// Initialize Core Services
	authService := services.NewAuthService(cfg.JWTConfig.Secret, authRepo)
	userService := services.NewUserService(userRepo)
	projectService := services.NewProjectService(projectRepo, projectProcessor)

	// Initialize Handlers
	userHandler := handlers.NewUserHandler(userService)
	authHandler := handlers.NewAuthHandler(userService, authService)
	projectHandler := handlers.NewProjectHandler(projectService)
	searchHandler := handlers.NewSearchHandler(queryService, summaryService)
	summaryHandler := handlers.NewSummaryHandler(summaryService)

	return &application{
		config:         cfg,
		authService:    authService,
		userService:    userService,
		projectService: projectService,
		queryService:   queryService,
		authHandler:    authHandler,
		userHandler:    userHandler,
		projectHandler: projectHandler,
		searchHandler:  searchHandler,
		summaryService: summaryService,
		summaryHandler: summaryHandler,
	}, nil
}
