package repositories

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
	SearchProjects(keyword string) ([]*models.Project, error)
	CreateProject(project *models.Project) error
	UpdateProject(project *models.Project) error
	DeleteProject(id string) error
}

type projectRepository struct {
	db *pgxpool.Pool
}

// NewProjectRepository creates a new PostgreSQL project repository
func NewProjectRepository(db *pgxpool.Pool) ProjectRepository {
	return &projectRepository{db: db}
}

func (r *projectRepository) GetProject(id string) (*models.Project, error) {
	query := `
    SELECT p.project_id, p.title, p.description, p.introduction, p.owner_id, 
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

	// Unmarshal JSONB to Go struct
	if err := json.Unmarshal([]byte(sectionsJSON), &p.Sections); err != nil {
		return nil, fmt.Errorf("error unmarshalling project sections: %v", err)
	}

	return &p, nil
}

func (r *projectRepository) SearchProjects(keyword string) ([]*models.Project, error) {
	query := `
    SELECT p.project_id, p.title, p.description, p.introduction, p.owner_id, 
           p.created_at, p.updated_at, p.project_sections::TEXT
    FROM projects p
    WHERE p.title ILIKE $1 OR p.description ILIKE $1
    ORDER BY p.created_at DESC
    `

	rows, err := r.db.Query(context.Background(), query, "%"+keyword+"%")
	if err != nil {
		return nil, fmt.Errorf("error searching projects: %v", err)
	}
	defer rows.Close()

	var projects []*models.Project
	for rows.Next() {
		var p models.Project
		var sectionsJSON string

		err := rows.Scan(
			&p.ProjectID,
			&p.Title,
			&p.Description,
			&p.Introduction,
			&p.OwnerID,
			&p.CreatedAt,
			&p.UpdatedAt,
			&sectionsJSON,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning project row: %v", err)
		}

		// Unmarshal JSONB to Go struct
		if err := json.Unmarshal([]byte(sectionsJSON), &p.Sections); err != nil {
			return nil, fmt.Errorf("error unmarshalling project sections: %v", err)
		}

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
	fmt.Println(p.Sections)
	sectionsJSON, _ := json.Marshal(p.Sections)
	fmt.Println(sectionsJSON)

	err := r.db.QueryRow(
		context.Background(),
		query,
        p.ProjectID,
        p.OwnerID,
		p.Title,
		p.Description,
        p.Status,
		p.Introduction,
		string(sectionsJSON),
		time.Now(),
		time.Now(),
	).Scan(&p.ProjectID)

	if err != nil {
		return fmt.Errorf("error creating project: %v", err)
	}
	return nil
}

func (r *projectRepository) UpdateProject(p *models.Project) error {
	query := `
    UPDATE projects
    SET title = $1, description = $2, introduction = $3, 
        updated_at = $4, project_sections = $5
    WHERE project_id = $6
    `

	sectionsJSON, _ := json.Marshal(p.Sections)
	now := time.Now()

	_, err := r.db.Exec(
		context.Background(),
		query,
		p.Title,
		p.Description,
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
