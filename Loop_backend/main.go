package main

import (
	"fmt"
	"log"
	"net/http"

	"Loop_backend/config"
	"Loop_backend/internal/handlers"
	"Loop_backend/internal/repositories"
	"Loop_backend/internal/services"
<<<<<<< HEAD
	tagservices "Loop_backend/internal/services/tags"
	neo4j "Loop_backend/platform/database/neo4j"
=======
>>>>>>> 4a2f436bed91636c5c2e3782993f5ab211ecfca7
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
<<<<<<< HEAD
	defer neo4j.Close()

	// Start Server
	serverAddr := fmt.Sprintf("%s:%d", cfg.ServerConfig.Host, cfg.ServerConfig.Port)
=======

	// Start Server
	serverAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
>>>>>>> 4a2f436bed91636c5c2e3782993f5ab211ecfca7
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
<<<<<<< HEAD
	if err := postgres.InitDB(&cfg.RelationalDatabaseConfig); err != nil {
		return nil, fmt.Errorf("failed to initialize database: %v", err)
	}
	if err := neo4j.InitGraph(&cfg.Neo4jConfig); err != nil {
		return nil, fmt.Errorf("failed to initialize graph database: %v", err)
=======
	// Initialize Database and Run Migrations
	dbConfig := &postgres.Config{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		Name:     cfg.Database.Name,
	}
	if err := postgres.InitDB(dbConfig); err != nil {
		return nil, fmt.Errorf("failed to initialize database: %v", err)
>>>>>>> 4a2f436bed91636c5c2e3782993f5ab211ecfca7
	}

	// Get Database Instance
	db := postgres.GetDB()
<<<<<<< HEAD
	graphDB := neo4j.GetDriver()
=======
>>>>>>> 4a2f436bed91636c5c2e3782993f5ab211ecfca7

	// Initialize Repositories
	userRepo := repositories.NewUserRepository(db)
	projectRepo := repositories.NewProjectRepository(db)
	authRepo := repositories.NewAuthRepository(db)
<<<<<<< HEAD
	graphRepo := repositories.NewGraphRepository(graphDB)

	// Initialize Services
	authService := services.NewAuthService(cfg.JWTConfig.Secret, authRepo)
	userService := services.NewUserService(userRepo)
	textProcessor := tagservices.NewTextProcessor()
	tagGenerationService := tagservices.NewTagGenerationService(textProcessor)
	projectService := services.NewProjectService(projectRepo, graphRepo, tagGenerationService)
=======

	// Initialize Graph Repository
	graphRepo, err := repositories.NewGraphRepository(
		cfg.Neo4j.URI,
		cfg.Neo4j.Username,
		cfg.Neo4j.Password,
	)
	if err != nil {
		log.Printf("Warning: Failed to initialize graph repository: %v", err)
		// Continue without graph DB for now
		graphRepo = nil
	}

	// Initialize Services
	authService := services.NewAuthService(cfg.JWT.Secret, authRepo)
	userService := services.NewUserService(userRepo)
	projectService := services.NewProjectService(projectRepo, graphRepo)
>>>>>>> 4a2f436bed91636c5c2e3782993f5ab211ecfca7

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
