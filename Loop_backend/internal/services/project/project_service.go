package project

import (
    "Loop_backend/internal/models"
    "Loop_backend/internal/repositories/project"
    "fmt"
    "github.com/google/uuid"
    "github.com/rs/zerolog"
    "github.com/rs/zerolog/log"
)

// ProjectService defines the interface for project operations
type ProjectService interface {
    CreateProject(project *models.Project) error
    GetProject(id uuid.UUID) (*models.Project, error)
    UpdateProject(project *models.Project) error
    DeleteProject(id uuid.UUID) error
    ListProjects() ([]*models.ProjectInfo, error)
    GetUserProjects(userID uuid.UUID) ([]*models.Project, error)
    GenerateProjectSummary(projectID string) (map[string]interface{}, error)
}

// DefaultProjectService implements ProjectService
type DefaultProjectService struct {
    projectRepo      project.ProjectRepository
    projectProcessor ProjectProcessor
    logger           zerolog.Logger
}

// NewProjectService creates a new project service instance
func NewProjectService(projectRepo project.ProjectRepository, projectProcessor ProjectProcessor) ProjectService {
    if projectRepo == nil {
        panic("projectRepo cannot be nil")
    }
    return &DefaultProjectService{
        projectRepo:      projectRepo,
        projectProcessor: projectProcessor,
        logger:           log.With().Str("component", "ProjectService").Logger(),
    }
}

// CreateProject creates a new project
func (s *DefaultProjectService) CreateProject(project *models.Project) error {
    s.logger.Info().Str("projectID", project.ProjectID.String()).Msg("Creating new project")
    
    // Create project in database
    if err := s.projectRepo.CreateProject(project); err != nil {
        s.logger.Error().Err(err).Str("projectID", project.ProjectID.String()).Msg("Failed to create project")
        return fmt.Errorf("failed to create project: %w", err)
    }

    // Process project with AI if processor is available
    if s.projectProcessor != nil {
        s.logger.Info().Str("projectID", project.ProjectID.String()).Msg("Analyzing new project")
        if err := s.projectProcessor.AnalyzeNewProject(project); err != nil {
            s.logger.Error().Err(err).Str("projectID", project.ProjectID.String()).Msg("Failed to analyze project")
            return fmt.Errorf("failed to analyze project: %w", err)
        }
    }
    
    s.logger.Info().Str("projectID", project.ProjectID.String()).Msg("Project created successfully")

    return nil
}

// GetProject retrieves a project by ID
func (s *DefaultProjectService) GetProject(id uuid.UUID) (*models.Project, error) {
    s.logger.Debug().Str("projectID", id.String()).Msg("Getting project")
    
    // Convert UUID to string for repository
    project, err := s.projectRepo.GetProject(id.String())
    if err != nil {
        s.logger.Error().Err(err).Str("projectID", id.String()).Msg("Failed to get project")
        return nil, fmt.Errorf("failed to get project: %w", err)
    }
    
    s.logger.Debug().Str("projectID", id.String()).Msg("Project retrieved successfully")
    return project, nil
}

// UpdateProject updates an existing project
func (s *DefaultProjectService) UpdateProject(project *models.Project) error {
    if err := s.projectRepo.UpdateProject(project); err != nil {
        return fmt.Errorf("failed to update project: %w", err)
    }

    if s.projectProcessor != nil {
        if err := s.projectProcessor.AnalyzeNewProject(project); err != nil {
            return fmt.Errorf("failed to re-analyze project: %w", err)
        }
    }

    return nil
}

// DeleteProject deletes a project by ID
func (s *DefaultProjectService) DeleteProject(id uuid.UUID) error {
    // Convert UUID to string for repository
    if err := s.projectRepo.DeleteProject(id.String()); err != nil {
        return fmt.Errorf("failed to delete project: %w", err)
    }
    return nil
}

// ListProjects retrieves all projects
func (s *DefaultProjectService) ListProjects() ([]*models.ProjectInfo, error) {
    // Use search with empty keyword to get all projects
    projectInfos, err := s.projectRepo.SearchProjects("")
    if err != nil {
        return nil, fmt.Errorf("failed to list projects: %w", err)
    }

    return projectInfos, nil
}

// GetUserProjects retrieves all projects for a user
func (s *DefaultProjectService) GetUserProjects(userID uuid.UUID) ([]*models.Project, error) {
    s.logger.Debug().Str("userID", userID.String()).Msg("Getting user projects")
    
    // Get projects by owner ID directly
    userProjects, err := s.projectRepo.GetProjectsByOwnerID(userID)
    if err != nil {
        s.logger.Error().Err(err).Str("userID", userID.String()).Msg("Failed to get user projects")
        return nil, fmt.Errorf("failed to get user projects: %w", err)
    }

    s.logger.Debug().Str("userID", userID.String()).Int("count", len(userProjects)).Msg("User projects retrieved successfully")
    return userProjects, nil
}

// GenerateProjectSummary generates a summary for a project
func (s *DefaultProjectService) GenerateProjectSummary(projectID string) (map[string]interface{}, error) {
    s.logger.Info().Str("projectID", projectID).Msg("Generating project summary")
    
    // Parse project ID
    id, err := uuid.Parse(projectID)
    if err != nil {
        s.logger.Error().Err(err).Str("projectID", projectID).Msg("Invalid project ID")
        return nil, fmt.Errorf("invalid project ID: %w", err)
    }

    // Get project
    project, err := s.GetProject(id)
    if err != nil {
        s.logger.Error().Err(err).Str("projectID", projectID).Msg("Failed to get project")
        return nil, fmt.Errorf("failed to get project: %w", err)
    }

    // Generate summary directly using the processor
    var summary string
    if s.projectProcessor != nil {
        s.logger.Debug().Str("projectID", projectID).Msg("Using processor to generate summary")
        summary, err = s.projectProcessor.GenerateProjectSummary(project)
        if err != nil {
            s.logger.Error().Err(err).Str("projectID", projectID).Msg("Failed to generate project summary")
            return nil, fmt.Errorf("failed to generate project summary: %w", err)
        }
    } else {
        s.logger.Debug().Str("projectID", projectID).Msg("No processor available, using project description as summary")
        summary = project.Description
    }

    // Return summary
    s.logger.Info().Str("projectID", projectID).Msg("Project summary generated successfully")
    return map[string]interface{}{
        "projectId":  project.ProjectID,
        "title":      project.Title,
        "summary":    summary,
        "sections":   project.Sections,
        "createdAt":  project.CreatedAt,
        "updatedAt":  project.UpdatedAt,
    }, nil
}
