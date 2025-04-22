package project

import (
	"Loop_backend/internal/dto"
	"Loop_backend/internal/models"
	"Loop_backend/internal/repositories/project"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// ProjectService handles project-related operations
type ProjectService interface {
	CreateProject(project *dto.CreateProjectRequest, ownerId uuid.UUID) error
	GetProject(id string) (*models.Project, error)
	UpdateProject(project *models.Project) error
	SearchProjects(keyword string) ([]*models.ProjectInfo, error)
	DeleteProject(id string) error
}

type projectService struct {
	projectRepo project.ProjectRepository
	// tagService   TagService
	// graphService models.KnowledgeGraphService
	analyzer models.ProjectProcessor
}

// NewProjectService creates a new project service instance
func NewProjectService(
	projectRepo project.ProjectRepository,
	analyzer models.ProjectProcessor,
) ProjectService {
	return &projectService{
		projectRepo: projectRepo,
		analyzer:    analyzer,
	}
}

func (s *projectService) SearchProjects(keyword string) ([]*models.ProjectInfo, error) {
	projects, err := s.projectRepo.SearchProjects(keyword)
	if err != nil {
		return nil, err
	}

	return projects, nil
}

// CreateProject creates a new project and triggers LLM analysis
func (s *projectService) CreateProject(projectRequest *dto.CreateProjectRequest, ownerId uuid.UUID) error {
	now := time.Now()
	projectID := uuid.New()

	// Create new project with validated data
	project := &models.Project{
		ProjectInfo: models.ProjectInfo{
			ProjectID:    projectID,
			OwnerID:      ownerId,
			Title:        projectRequest.Title,
			Description:  projectRequest.Description,
			Status:       models.Status(projectRequest.Status),
			Introduction: projectRequest.Introduction,
			Tags:         projectRequest.Tags,
			CreatedAt:    now,
			UpdatedAt:    now,
		},
		Sections: projectRequest.Sections,
	}

	// Create project in repository
	if err := s.projectRepo.CreateProject(project); err != nil {
		return fmt.Errorf("failed to create project: %w", err)
	}

	if err := s.analyzer.AnalyzeNewProject(project); err != nil {
		return fmt.Errorf("failed to analyze project %s: %v", project.ProjectID, err)
	}

	// // Create graph representation if tags are present
	// if len(project.Tags) > 0 {
	//     if err := s.graphService.StoreProjectGraph(projectID, &models.KnowledgeGraph{
	//         Entities: s.convertTagsToEntities(project.Tags),
	//         Relationships: []models.Relationship{},
	//     }); err != nil {
	//         log.Printf("Warning: failed to create project graph for %s: %v", project.ProjectID, err)
	//     }
	// }

	return nil
}

// GetProject retrieves a project by ID
func (s *projectService) GetProject(id string) (*models.Project, error) {
	return s.projectRepo.GetProject(id)
}

// UpdateProject updates an existing project
func (s *projectService) UpdateProject(project *models.Project) error {
	project.UpdatedAt = time.Now()
	return s.projectRepo.UpdateProject(project)
}

// DeleteProject deletes a project
func (s *projectService) DeleteProject(id string) error {
	return s.projectRepo.DeleteProject(id)
}
