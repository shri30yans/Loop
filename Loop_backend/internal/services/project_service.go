package services

import (
	"Loop_backend/internal/dto"
	"Loop_backend/internal/models"
	"Loop_backend/internal/repositories"
	tagservice "Loop_backend/internal/services/tags"
	"fmt"
)

type ProjectService interface {
	GetProject(project_id string) (*models.Project, error)
	SearchProjects(keyword string) ([]*models.Project, error)
	CreateProject(req dto.CreateProjectRequest) (*models.Project, error)
	DeleteProject(project_id string) error
}

type projectService struct {
	pgRepo           repositories.ProjectRepository
	graphRepo        repositories.GraphRepository
	tagGenerationSvc tagservice.TagGenerationService
}

func NewProjectService(pgRepo repositories.ProjectRepository, graphRepo repositories.GraphRepository, tagGenerationSvc tagservice.TagGenerationService) ProjectService {
	return &projectService{
		pgRepo:           pgRepo,
		graphRepo:        graphRepo,
		tagGenerationSvc: tagGenerationSvc,
	}
}

func (s *projectService) GetProject(project_id string) (*models.Project, error) {
	// Get project data from PostgreSQL
	project, err := s.pgRepo.GetProject(project_id)
	if err != nil {
		return nil, err
	}
	return project, nil
}

func (s *projectService) SearchProjects(keyword string) ([]*models.Project, error) {
	projects, err := s.pgRepo.SearchProjects(keyword)
	if err != nil {
		return nil, err
	}

	return projects, nil
}

func (s *projectService) CreateProject(req dto.CreateProjectRequest) (*models.Project, error) {
	newProject, err := models.NewProject(
		req.OwnerID,
		req.Title,
		req.Description,
		req.Status,
		req.Introduction,
		req.Tags,
		req.Sections,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create project model: %v", err)
	}

	// Generate tags automatically
	generatedTags, err := s.tagGenerationSvc.GenerateProjectTags(newProject)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tags: %v", err)
	}

	// Merge user-provided and generated tags
	newProject.Tags = s.tagGenerationSvc.MergeTags(req.Tags, generatedTags)

	if err := s.pgRepo.CreateProject(newProject); err != nil {
		return nil, fmt.Errorf("failed to save project to PostgreSQL: %v", err)
	}

	// Create project with user and tags in Neo4j
	if err := s.graphRepo.CreateProjectWithUserAndTags(newProject, newProject.Tags); err != nil {
		return nil, fmt.Errorf("failed to create project in GraphDB: %v", err)
	}

	// // Set tags in the returned project object
	// if len(req.Tags) > 0 {
	// 	newProject.Tags = req.Tags
	// } else {
	// 	newProject.Tags = []string{} // Empty array for consistency
	// }

	return newProject, nil
}

func (s *projectService) DeleteProject(project_id string) error {
	if err := s.pgRepo.DeleteProject(project_id); err != nil {
		return err
	}

	if err := s.graphRepo.DeleteProjectNode(project_id); err != nil {
		fmt.Printf("Warning: Failed to delete project from graph DB: %v\n", err)
	}

	return nil
}
