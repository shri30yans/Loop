package main

import (
    "context"
    "fmt"
    "net/http"
    "os"

    "Loop_backend/config"
    "Loop_backend/internal/handlers"
    "Loop_backend/internal/ai/providers"
    authRepo "Loop_backend/internal/repositories/auth"
    projectRepo "Loop_backend/internal/repositories/project"
    userRepo "Loop_backend/internal/repositories/user"
    vectorRepo "Loop_backend/internal/repositories/vector"
    graphRepo "Loop_backend/internal/repositories/graph"
    "Loop_backend/internal/services/auth"
    "Loop_backend/internal/services/project"
    "Loop_backend/internal/services/search"
    "Loop_backend/internal/services/user"
    "Loop_backend/internal/services/graph"
    "Loop_backend/internal/services/embedding"
    "Loop_backend/internal/services/graph/response_parser"

    "github.com/joho/godotenv"
    "github.com/rs/cors"
    "github.com/rs/zerolog"
)

type application struct {
	config         *config.Config
	authService    auth.AuthService
	userService    user.UserService
	projectService project.ProjectService
	authHandler    *handlers.AuthHandler
	userHandler    *handlers.UserHandler
	projectHandler *handlers.ProjectHandler
	searchHandler  *handlers.SearchHandler
	toolsHandler   *handlers.ToolsHandler
}

func main() {
    // Initialize structured logging
    logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
    
    // Load environment variables
    if err := godotenv.Load(); err != nil {
        logger.Info().Msg("No .env file found")
    }

    // Load Configuration
    ctx := context.Background()
    cfg, err := config.LoadConfig(ctx)
    if err != nil {
        logger.Fatal().Err(err).Msg("Failed to load config")
    }

    // Initialize all components
    app, err := initializeApp(cfg, &logger)
    if err != nil {
        logger.Fatal().Err(err).Msg("Failed to initialize application")
    }

    // Get auth middleware once app is initialized
    authMiddleware := app.authService.AuthenticateRequest

    // Initialize Router and RouteRegister with auth middleware
    routeRegister := handlers.NewRouteRegister(authMiddleware)

    // Register routes using the initialized app
    app.authHandler.RegisterRoutes(routeRegister)
    app.userHandler.RegisterRoutes(routeRegister)
    app.projectHandler.RegisterRoutes(routeRegister)
    app.toolsHandler.RegisterRoutes(routeRegister)

    // Get the underlying router
    router := routeRegister.Router()

    // Setup CORS
    c := cors.New(cors.Options{
        AllowedOrigins:   []string{"http://localhost:3000"},
        AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
        AllowedHeaders:   []string{"Authorization", "Content-Type"},
        AllowCredentials: true,
    })

    // Wrap the router with the CORS handler
    handler := c.Handler(router)

    // Start Server
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    serverAddr := fmt.Sprintf(":%s", port)
    logger.Info().Str("address", serverAddr).Msg("Server running")

    server := &http.Server{
        Addr:    serverAddr,
        Handler: handler,
    }

    if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
        logger.Fatal().Err(err).Msg("Server failed")
    }
}

func initializeApp(cfg *config.Config, logger *zerolog.Logger) (*application, error) {
    // Get Database Instance
    db := cfg.PostgresPool
    if db == nil {
        return nil, fmt.Errorf("database connection pool is nil")
    }
    
    logger.Info().Msg("Database connection established")

    // Initialize Repositories
    uRepo := userRepo.NewUserRepository(db)
    pRepo := projectRepo.NewProjectRepository(db)
    if pRepo == nil {
        return nil, fmt.Errorf("failed to create project repository")
    }
    aRepo := authRepo.NewAuthRepository(db)
    
    // Create a new connection for vector repository since it needs v5
    dbv5, err := vectorRepo.InitDB(cfg.RelationalDatabaseConfig.ConnectionString)
    if err != nil {
        return nil, fmt.Errorf("failed to initialize vector db: %v", err)
    }
    vRepo := vectorRepo.NewVectorRepository(dbv5)

    // Initialize Graph Repository
    gRepo := graphRepo.NewRepository(cfg.Neo4jDriver)

    // Initialize AI Provider
    aiProvider, err := providers.NewProvider(&cfg.AIConfig)
    if err != nil {
        return nil, fmt.Errorf("failed to initialize AI provider: %v", err)
    }

    // Initialize Graph Service
    graphService := graph.NewService(gRepo, aiProvider)

    // Initialize Python Embedding Client
    embeddingClient := embedding.NewPythonEmbeddingClient(cfg.AIConfig.PythonEmbeddingURL)

    // Initialize Response Parser
    responseParser := response_parser.NewResponseParser()

    // Initialize Project Processor with the graph service interface
    projectProcessor := project.NewProjectProcessor(
        graphService,
        vRepo,
        aiProvider,
        embeddingClient,
        responseParser,
    )
    
    logger.Info().Msg("Project processor initialized")

    // Initialize Core Services
    userService := user.NewUserService(uRepo)
    projectService := project.NewProjectService(pRepo, projectProcessor)
    
    logger.Info().Msg("Core services initialized")
    authService := auth.NewAuthService(cfg.JWTConfig.Secret, aRepo)

    // Get auth middleware
    authMiddleware := authService.AuthenticateRequest

    // Initialize Vector Search Service
    vectorSearchService := search.NewVectorSearchService(vRepo, embeddingClient)
    
    // Initialize Search Service with only vector search capabilities
    // Pass nil for graphService to disable graph-based search
    queryService := search.NewHybridSearchService(nil, vectorSearchService)

    // Initialize Handlers
    authHandler := handlers.NewAuthHandler(userService, authService)
    userHandler := handlers.NewUserHandler(userService, authMiddleware)
    projectHandler := handlers.NewProjectHandler(projectService, authMiddleware)
    searchHandler := handlers.NewSearchHandler(queryService)
    
    // Get API token from config
    apiToken := os.Getenv("GO_BACKEND_API_TOKEN")
    toolsHandler := handlers.NewToolsHandler(queryService, apiToken)

    return &application{
    	config:         cfg,
    	authService:    authService,
    	userService:    userService,
    	projectService: projectService,
    	authHandler:    authHandler,
    	userHandler:    userHandler,
    	projectHandler: projectHandler,
    	searchHandler:  searchHandler,
    	toolsHandler:   toolsHandler,
    }, nil
}
