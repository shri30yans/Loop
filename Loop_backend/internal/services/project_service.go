package services

import (
	"Loop_backend/internal/dto"
	"Loop_backend/internal/models"
	"Loop_backend/internal/repositories"
<<<<<<< HEAD
	tagservice "Loop_backend/internal/services/tags"
=======
>>>>>>> 4a2f436bed91636c5c2e3782993f5ab211ecfca7
	"fmt"
)

type ProjectService interface {
	GetProject(project_id string) (*models.Project, error)
<<<<<<< HEAD
	SearchProjects(keyword string) ([]*models.Project, error)
=======
	GetProjects() ([]*models.Project, error)
	SearchProjects(keyword string) ([]*models.Project, int, error)
>>>>>>> 4a2f436bed91636c5c2e3782993f5ab211ecfca7
	CreateProject(req dto.CreateProjectRequest) (*models.Project, error)
	DeleteProject(project_id string) error
}

type projectService struct {
<<<<<<< HEAD
	pgRepo           repositories.ProjectRepository
	graphRepo        repositories.GraphRepository
	tagGenerationSvc tagservice.TagGenerationService
}

func NewProjectService(pgRepo repositories.ProjectRepository, graphRepo repositories.GraphRepository, tagGenerationSvc tagservice.TagGenerationService) ProjectService {
	return &projectService{
		pgRepo:           pgRepo,
		graphRepo:        graphRepo,
		tagGenerationSvc: tagGenerationSvc,
=======
	pgRepo    repositories.ProjectRepository
	graphRepo repositories.GraphRepository
}

func NewProjectService(pgRepo repositories.ProjectRepository, graphRepo repositories.GraphRepository) ProjectService {
	return &projectService{
		pgRepo:    pgRepo,
		graphRepo: graphRepo,
>>>>>>> 4a2f436bed91636c5c2e3782993f5ab211ecfca7
	}
}

func (s *projectService) GetProject(project_id string) (*models.Project, error) {
	// Get project data from PostgreSQL
	project, err := s.pgRepo.GetProject(project_id)
	if err != nil {
		return nil, err
	}
<<<<<<< HEAD
	return project, nil
}

func (s *projectService) SearchProjects(keyword string) ([]*models.Project, error) {
	projects, err := s.pgRepo.SearchProjects(keyword)
=======

	// Get tags from Neo4j
	tags, err := s.graphRepo.GetProjectTags(project_id)
	if err != nil {
		// Log error but don't fail the request
		fmt.Printf("Warning: Failed to get project tags from graph DB: %v\n", err)
		project.Tags = []string{} // Use empty array as fallback
	} else {
		project.Tags = tags
	}

	return project, nil
}

func (s *projectService) GetProjects() ([]*models.Project, error) {
	// Get projects from PostgreSQL
	projects, err := s.pgRepo.GetProjects()
>>>>>>> 4a2f436bed91636c5c2e3782993f5ab211ecfca7
	if err != nil {
		return nil, err
	}

<<<<<<< HEAD
	return projects, nil
}

func (s *projectService) CreateProject(req dto.CreateProjectRequest) (*models.Project, error) {
=======
	// Get project IDs
	projectIDs := make([]string, len(projects))
	for i, project := range projects {
		projectIDs[i] = project.ProjectID
	}

	// Get tags for all projects from Neo4j
	projectTags, err := s.graphRepo.GetProjectsWithTags(projectIDs)
	if err != nil {
		// Log error but don't fail the request
		fmt.Printf("Warning: Failed to get project tags from graph DB: %v\n", err)
	} else {
		// Associate tags with projects
		for _, project := range projects {
			if tags, ok := projectTags[project.ProjectID]; ok {
				project.Tags = tags
			} else {
				project.Tags = []string{} // Use empty array for projects without tags
			}
		}
	}

	return projects, nil
}

func (s *projectService) SearchProjects(keyword string) ([]*models.Project, int, error) {
	// Search projects in PostgreSQL
	projects, err := s.pgRepo.SearchProjects(keyword)
	if err != nil {
		return nil, 0, err
	}

	// Get project IDs
	projectIDs := make([]string, len(projects))
	for i, project := range projects {
		projectIDs[i] = project.ProjectID
	}

	// Get tags for all projects from Neo4j
	projectTags, err := s.graphRepo.GetProjectsWithTags(projectIDs)
	if err != nil {
		// Log error but don't fail the request
		fmt.Printf("Warning: Failed to get project tags from graph DB: %v\n", err)
	} else {
		// Associate tags with projects
		for _, project := range projects {
			if tags, ok := projectTags[project.ProjectID]; ok {
				project.Tags = tags
			} else {
				project.Tags = []string{} // Use empty array for projects without tags
			}
		}
	}

	return projects, len(projects), nil
}

func (s *projectService) CreateProject(req dto.CreateProjectRequest) (*models.Project, error) {
	// Create project in PostgreSQL without tags
>>>>>>> 4a2f436bed91636c5c2e3782993f5ab211ecfca7
	newProject, err := models.NewProject(
		req.OwnerID,
		req.Title,
		req.Description,
		req.Status,
		req.Introduction,
<<<<<<< HEAD
		req.Tags,
=======
		[]string{}, // Empty tags for PostgreSQL
>>>>>>> 4a2f436bed91636c5c2e3782993f5ab211ecfca7
		req.Sections,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create project model: %v", err)
	}

<<<<<<< HEAD
	// Generate tags automatically
	generatedTags, err := s.tagGenerationSvc.GenerateProjectTags(newProject)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tags: %v", err)
	}

	// Merge user-provided and generated tags
	newProject.Tags = s.tagGenerationSvc.MergeTags(req.Tags, generatedTags)

=======
>>>>>>> 4a2f436bed91636c5c2e3782993f5ab211ecfca7
	if err := s.pgRepo.CreateProject(newProject); err != nil {
		return nil, fmt.Errorf("failed to save project to PostgreSQL: %v", err)
	}

<<<<<<< HEAD
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
=======
	// Create project node and user-project relationship in Neo4j
	if err := s.graphRepo.CreateProjectNode(newProject); err != nil {
		// Log error but don't fail the request
		fmt.Printf("Warning: Failed to create project node in graph DB: %v\n", err)
	}

	// Create project-tag relationships if there are any tags
	if len(req.Tags) > 0 {
		if err := s.graphRepo.CreateProjectTagRelations(newProject.ProjectID, req.Tags); err != nil {
			fmt.Printf("Warning: Failed to create project-tag relations in graph DB: %v\n", err)
		}
		// Set tags in the returned project object
		newProject.Tags = req.Tags
	} else {
		newProject.Tags = []string{} // Set empty array for consistency
	}
>>>>>>> 4a2f436bed91636c5c2e3782993f5ab211ecfca7

	return newProject, nil
}

func (s *projectService) DeleteProject(project_id string) error {
<<<<<<< HEAD
=======
	// First delete from PostgreSQL
>>>>>>> 4a2f436bed91636c5c2e3782993f5ab211ecfca7
	if err := s.pgRepo.DeleteProject(project_id); err != nil {
		return err
	}

<<<<<<< HEAD
=======
	// Delete from Neo4j (including all relationships)
>>>>>>> 4a2f436bed91636c5c2e3782993f5ab211ecfca7
	if err := s.graphRepo.DeleteProjectNode(project_id); err != nil {
		fmt.Printf("Warning: Failed to delete project from graph DB: %v\n", err)
	}

	return nil
}
