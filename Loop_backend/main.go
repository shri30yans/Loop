package main

import (
    "fmt"
    "log"
    "net/http"

    "github.com/rs/cors"
    "Loop_backend/config"
    "Loop_backend/internal/handlers"
    "Loop_backend/internal/repositories"
    "Loop_backend/internal/services"
    postgres "Loop_backend/platform/database/postgres"
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

    // Start Server
    serverAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
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
    }

    // Get Database Instance
    db := postgres.GetDB()

    // Initialize Repositories
    userRepo := repositories.NewUserRepository(db)
    projectRepo := repositories.NewProjectRepository(db)
    authRepo := repositories.NewAuthRepository(db)

    // Initialize Services
    authService := services.NewAuthService(cfg.JWT.Secret, authRepo)
    userService := services.NewUserService(userRepo)
    projectService := services.NewProjectService(projectRepo)

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
