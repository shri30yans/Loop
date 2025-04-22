package project

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"Loop_backend/internal/models"

	"github.com/jackc/pgx/v4/pgxpool"
)

type ProjectRepository interface {
	GetProject(id string) (*models.Project, error)
	SearchProjects(keyword string) ([]*models.ProjectInfo, error)
	CreateProject(project *models.Project) error
	UpdateProject(project *models.Project) error
	DeleteProject(id string) error
}

type projectRepository struct {
	db *pgxpool.Pool
}

func NewProjectRepository(db *pgxpool.Pool) ProjectRepository {
	return &projectRepository{db: db}
}

func (r *projectRepository) GetProject(id string) (*models.Project, error) {
	query := `
        SELECT p.project_id, p.title, p.description, p.status, p.introduction, p.owner_id, 
               p.created_at, p.updated_at, p.project_sections::TEXT
        FROM projects p
        WHERE p.project_id = $1
        `
	var p models.Project
	var sectionsJSON string

	err := r.db.QueryRow(context.Background(), query, id).Scan(
		&p.ProjectID,
		&p.Title,
		&p.Description,
		&p.Status,
		&p.Introduction,
		&p.OwnerID,
		&p.CreatedAt,
		&p.UpdatedAt,
		&sectionsJSON,
	)

	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, fmt.Errorf("project not found: %v", err)
		}
		return nil, fmt.Errorf("error finding project: %v", err)
	}

	// Initialize empty tags array
	p.Tags = []string{}

	// Unmarshal JSONB to Go struct
	if err := json.Unmarshal([]byte(sectionsJSON), &p.Sections); err != nil {
		return nil, fmt.Errorf("error unmarshalling project sections: %v", err)
	}

	return &p, nil
}

func (r *projectRepository) SearchProjects(keyword string) ([]*models.ProjectInfo, error) {
	query := `
        SELECT p.project_id, p.title, p.description, p.status, p.introduction, p.owner_id, 
               p.created_at, p.updated_at
        FROM projects p
        WHERE p.title ILIKE $1 OR p.description ILIKE $1
        ORDER BY p.created_at DESC
        `

	rows, err := r.db.Query(context.Background(), query, "%"+keyword+"%")
	if err != nil {
		return nil, fmt.Errorf("error searching projects: %v", err)
	}
	defer rows.Close()

	var projects []*models.ProjectInfo
	for rows.Next() {
		var p models.ProjectInfo

		err := rows.Scan(
			&p.ProjectID,
			&p.Title,
			&p.Description,
			&p.Status,
			&p.Introduction,
			&p.OwnerID,
			&p.CreatedAt,
			&p.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning project row: %v", err)
		}

		// Initialize empty tags array
		p.Tags = []string{}

		projects = append(projects, &p)
	}
	return projects, nil
}

func (r *projectRepository) CreateProject(p *models.Project) error {
	query := `
        INSERT INTO projects (project_id, owner_id, title, description, status, introduction, project_sections, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
        RETURNING project_id
        `

	sectionsJSON, err := json.Marshal(p.Sections)
	if err != nil {
		return fmt.Errorf("error marshalling sections: %v", err)
	}

	now := time.Now()

	err = r.db.QueryRow(
		context.Background(),
		query,
		p.ProjectID,
		p.OwnerID,
		p.Title,
		p.Description,
		p.Status,
		p.Introduction,
		string(sectionsJSON),
		now,
		now,
	).Scan(&p.ProjectID)

	if err != nil {
		return fmt.Errorf("error creating project: %v", err)
	}
	return nil
}

func (r *projectRepository) UpdateProject(p *models.Project) error {
	query := `
        UPDATE projects
        SET title = $1, description = $2, status = $3, introduction = $4, 
            updated_at = $5, project_sections = $6
        WHERE project_id = $7
        `

	sectionsJSON, err := json.Marshal(p.Sections)
	if err != nil {
		return fmt.Errorf("error marshalling sections: %v", err)
	}

	now := time.Now()

	_, err = r.db.Exec(
		context.Background(),
		query,
		p.Title,
		p.Description,
		p.Status,
		p.Introduction,
		now,
		string(sectionsJSON),
		p.ProjectID,
	)
	if err != nil {
		return fmt.Errorf("error updating project: %v", err)
	}

	p.UpdatedAt = now
	return nil
}

func (r *projectRepository) DeleteProject(id string) error {
	result, err := r.db.Exec(context.Background(), "DELETE FROM projects WHERE project_id = $1", id)
	if err != nil {
		return fmt.Errorf("error deleting project: %v", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("project not found")
	}

	return nil
}
