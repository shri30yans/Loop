package services

import (
	"Loop_backend/internal/models"
	"Loop_backend/internal/repositories"
	"Loop_backend/internal/dto"
	"fmt"
)


type ProjectService interface {
	GetProject(project_id string) (*models.Project, error)
	SearchProjects(keyword string) ([]*models.Project, int, error)
	CreateProject(req dto.CreateProjectRequest) (*models.Project, error)
	DeleteProject(project_id string) error
}

type projectService struct {
	repo repositories.ProjectRepository
}

// NewProjectService creates a new project service
func NewProjectService(repo repositories.ProjectRepository) ProjectService {
	return &projectService{repo: repo}
}

func (s *projectService) GetProject(project_id string) (*models.Project, error) {
	return s.repo.GetProject(project_id)
}


func (s *projectService) SearchProjects(keyword string) ([]*models.Project, int, error) {
	
	projects, err := s.repo.SearchProjects(keyword)
	if err != nil {
		return nil, 0, nil
	}
	return projects, len(projects), nil
}

func (s *projectService) CreateProject(req dto.CreateProjectRequest) (*models.Project, error) {
	fmt.Println(req.Sections)
	newProject, _ := models.NewProject(
		req.OwnerID,
		req.Title,
		req.Description,
		req.Status,
		req.Introduction,
		req.Tags,
		req.Sections,
	)

	if err := s.repo.CreateProject(newProject); err != nil {
		return nil, err
	}

	return newProject, nil
}

func (s *projectService) DeleteProject(project_id string) error {
	return s.repo.DeleteProject(project_id)
}
