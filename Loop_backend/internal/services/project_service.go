package services

import (
	"Loop_backend/internal/models"
	"Loop_backend/internal/repositories"
	"errors"
)

// Project DTO types
type CreateProjectRequest struct {
	Title        string           `json:"title"`
	Description  string           `json:"description"`
	Introduction string           `json:"introduction"`
	OwnerID      string              `json:"owner_id"`
	Tags         []string         `json:"tags"`
	Sections     []SectionRequest `json:"sections"`
}

type SectionRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type ProjectService interface {
	GetProject(project_id string) (*models.Project, error)
	GetUserProjects(ownerID string) ([]*models.Project, error)
	SearchProjects(keyword string) ([]*models.Project, int, error)
	CreateProject(req CreateProjectRequest) (*models.Project, error)
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
	return s.repo.FindByID(project_id)
}

func (s *projectService) GetUserProjects(ownerID string) ([]*models.Project, error) {
	return s.repo.FindByOwner(ownerID)
}

func (s *projectService) SearchProjects(keyword string) ([]*models.Project, int, error) {
	if keyword == "" {
		return nil, 0, errors.New("empty search keyword")
	}
	return s.repo.Search(keyword)
}

func (s *projectService) CreateProject(req CreateProjectRequest) (*models.Project, error) {
	if req.Title == "" {
		return nil, models.ErrInvalidTitle
	}

	newProject := &models.Project{
		OwnerID:      string(req.OwnerID),
		Title:        req.Title,
		Description:  req.Description,
		Introduction: req.Introduction,
		Tags:         req.Tags,
		Sections:     make([]models.Section, len(req.Sections)),
	}

	// Convert SectionRequest to models.Section
	for i, s := range req.Sections {
		section, err := models.NewSection(s.Title, s.Content)
		if err != nil {
			return nil, err
		}
		newProject.Sections[i] = *section
	}

	if err := s.repo.CreateProject(newProject); err != nil {
		return nil, err
	}

	return newProject, nil
}

func (s *projectService) DeleteProject(project_id string) error {
	return s.repo.Delete(project_id)
}
