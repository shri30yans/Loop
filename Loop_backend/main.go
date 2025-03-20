package main

import (
	"fmt"
	"log"
	"net/http"

	"Loop_backend/config"
	"Loop_backend/internal/handlers"
	"Loop_backend/internal/repositories"
	"Loop_backend/internal/services"
	tagservices "Loop_backend/internal/services/tags"
	neo4j "Loop_backend/platform/database/neo4j"
	postgres "Loop_backend/platform/database/postgres"

	"github.com/rs/cors"
)

type application struct {
	config         *config.Config
	authService    services.AuthService
	userService    services.UserService
	projectService services.ProjectService
	authHandler    *handlers.AuthHandler
	userHandler    *handlers.UserHandler
	projectHandler *handlers.ProjectHandler
}

func main() {
	// Load Configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	app, err := initializeApp(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	// Initialize Router
	mux := app.routes()

	// Setup CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})

	// Wrap the router with the CORS handler
	handler := c.Handler(mux)

	// Setup graceful shutdown
	defer postgres.Close()
	defer neo4j.Close()

	// Start Server
	serverAddr := fmt.Sprintf("%s:%d", cfg.ServerConfig.Host, cfg.ServerConfig.Port)
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
	if err := neo4j.InitGraph(&cfg.Neo4jConfig); err != nil {
		return nil, fmt.Errorf("failed to initialize graph database: %v", err)
	}

	// Get Database Instance
	db := postgres.GetDB()
	graphDB := neo4j.GetDriver()

	// Initialize Repositories
	userRepo := repositories.NewUserRepository(db)
	projectRepo := repositories.NewProjectRepository(db)
	authRepo := repositories.NewAuthRepository(db)
	graphRepo := repositories.NewGraphRepository(graphDB)

	// Initialize Services
	authService := services.NewAuthService(cfg.JWTConfig.Secret, authRepo)
	userService := services.NewUserService(userRepo)
	textProcessor := tagservices.NewTextProcessor()
	tagGenerationService := tagservices.NewTagGenerationService(textProcessor)
	projectService := services.NewProjectService(projectRepo, graphRepo, tagGenerationService)

	// Initialize Handlers
	userHandler := handlers.NewUserHandler(userService)
	authHandler := handlers.NewAuthHandler(userService, authService)
	projectHandler := handlers.NewProjectHandler(projectService)

	return &application{
		config:         cfg,
		authService:    authService,
		userService:    userService,
		projectService: projectService,
		authHandler:    authHandler,
		userHandler:    userHandler,
		projectHandler: projectHandler,
	}, nil
}

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	// Create a new RouteRegister that handles auth middleware
	routeRegister := handlers.NewRouteRegister(mux, app.authService)

	// Register routes for all handlers
	app.authHandler.RegisterRoutes(routeRegister)
	app.userHandler.RegisterRoutes(routeRegister)
	app.projectHandler.RegisterRoutes(routeRegister)

	return mux
}
