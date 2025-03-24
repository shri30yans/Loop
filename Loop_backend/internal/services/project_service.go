package services

import (
    "fmt"

    "Loop_backend/internal/dto"
    "Loop_backend/internal/models"
    "Loop_backend/internal/repositories"
    "Loop_backend/internal/services/tags"
)

type ProjectService interface {
    GetProject(project_id string) (*models.Project, error)
    SearchProjects(keyword string) ([]*models.Project, error)
    CreateProject(req dto.CreateProjectRequest) (*models.Project, error)
    DeleteProject(project_id string) error
}

type projectService struct {
    pgRepo              repositories.ProjectRepository
    graphRepo           repositories.GraphRepository
    tagGenerationSvc    tags.TagGenerationService
    entityProcessingSvc EntityProcessingService
}

func NewProjectService(
    pgRepo repositories.ProjectRepository,
    graphRepo repositories.GraphRepository,
    tagService tags.TagGenerationService,
    entityProcessingSvc EntityProcessingService,
) ProjectService {
    return &projectService{
        pgRepo:              pgRepo,
        graphRepo:           graphRepo,
        tagGenerationSvc:    tagService,
        entityProcessingSvc: entityProcessingSvc,
    }
}

func (s *projectService) GetProject(project_id string) (*models.Project, error) {
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

func (s *projectService) CreateTags(project *models.Project) ([]models.Tag, error) {
    // Generate tags automatically
    generatedTags, err := s.tagGenerationSvc.GenerateProjectTags(project)
    if err != nil {
        return nil, fmt.Errorf("failed to generate tags: %v", err)
    }
    return generatedTags, nil
}

// extractTagNames helper function to get tag names from Tag structs
func extractTagNames(tags []models.Tag) []string {
    var names []string
    for _, tag := range tags {
        names = append(names, tag.Name)
    }
    return names
}

func (s *projectService) CreateProject(req dto.CreateProjectRequest) (*models.Project, error) {
    // Create new project instance
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
        return nil, err
    }

    // Generate tags
    generatedTags, err := s.CreateTags(newProject)
    if err != nil {
        return nil, fmt.Errorf("failed to generate tags: %v", err)
    }

    // Save project to PostgreSQL
    if err := s.pgRepo.CreateProject(newProject); err != nil {
        return nil, fmt.Errorf("failed to save project to PostgreSQL: %v", err)
    }

    // Create project with user and tags in Neo4j
    tagNames := extractTagNames(generatedTags)
    if err := s.graphRepo.CreateProjectWithUserAndTags(newProject, tagNames); err != nil {
        return nil, fmt.Errorf("failed to create project in GraphDB: %v", err)
    }

    // Process project entities and relationships
    if err := s.entityProcessingSvc.ProcessProjectEntities(newProject); err != nil {
        return nil, fmt.Errorf("failed to process project entities: %v", err)
    }

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
