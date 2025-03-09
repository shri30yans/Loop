package repositories

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"Loop_backend/internal/models"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
)

type ProjectRepository interface {
	FindByID(id string) (*models.Project, error)
	FindByOwner(ownerID string) ([]*models.Project, error)
	Search(keyword string) ([]*models.Project, int, error)
	CreateProject(project *models.Project) error
	Update(project *models.Project) error
	Delete(id string) error
}

type projectRepository struct {
	db *pgxpool.Pool
}

// NewProjectRepository creates a new PostgreSQL project repository
func NewProjectRepository(db *pgxpool.Pool) ProjectRepository {
	return &projectRepository{db: db}
}

func (r *projectRepository) FindByID(id string) (*models.Project, error) {
	query := `
    SELECT p.project_id, p.title, p.description, p.introduction, p.owner_id, 
           p.created_at, p.updated_at, p.project_sections::TEXT
    FROM projects p
    WHERE p.project_id = $1
    `

	var p models.Project
	var sectionsJSON string

	err := r.db.QueryRow(context.Background(), query, id).Scan(
		&p.ID,
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

func (r *projectRepository) FindByOwner(ownerID string) ([]*models.Project, error) {
	query := `
    SELECT p.project_id, p.title, p.description, p.introduction, p.owner_id, 
           p.created_at, p.updated_at, p.project_sections::TEXT
    FROM projects p
    WHERE p.owner_id = $1
    ORDER BY p.created_at DESC
    `

	rows, err := r.db.Query(context.Background(), query, ownerID)
	if err != nil {
		return nil, fmt.Errorf("error querying projects: %v", err)
	}
	defer rows.Close()

	var projects []*models.Project
	for rows.Next() {
		var p models.Project
		var sectionsJSON string

		err := rows.Scan(
			&p.ID,
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

	if len(projects) == 0 {
		return nil, fmt.Errorf("no projects found for owner")
	}

	return projects, nil
}

func (r *projectRepository) Search(keyword string) ([]*models.Project, int, error) {
	query := `
    SELECT p.project_id, p.title, p.description, p.introduction, p.owner_id, 
           p.created_at, p.updated_at, p.project_sections::TEXT
    FROM projects p
    WHERE p.title ILIKE $1 OR p.description ILIKE $1
    ORDER BY p.created_at DESC
    `

	rows, err := r.db.Query(context.Background(), query, "%"+keyword+"%")
	if err != nil {
		return nil, 0, fmt.Errorf("error searching projects: %v", err)
	}
	defer rows.Close()

	var projects []*models.Project
	for rows.Next() {
		var p models.Project
		var sectionsJSON string

		err := rows.Scan(
			&p.ID,
			&p.Title,
			&p.Description,
			&p.Introduction,
			&p.OwnerID,
			&p.CreatedAt,
			&p.UpdatedAt,
			&sectionsJSON,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("error scanning project row: %v", err)
		}

		// Unmarshal JSONB to Go struct
		if err := json.Unmarshal([]byte(sectionsJSON), &p.Sections); err != nil {
			return nil, 0, fmt.Errorf("error unmarshalling project sections: %v", err)
		}

		projects = append(projects, &p)
	}

	var total int
	countQuery := `
    SELECT COUNT(*) FROM projects 
    WHERE title ILIKE $1 OR description ILIKE $1
    `
	err = r.db.QueryRow(context.Background(), countQuery, "%"+keyword+"%").Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("error counting projects: %v", err)
	}

	return projects, total, nil
}

func (r *projectRepository) CreateProject(p *models.Project) error {
	query := `
    INSERT INTO projects (title, description, introduction, owner_id, 
                          created_at, updated_at, project_sections)
    VALUES ($1, $2, $3, $4, $5, $5, $6)
    RETURNING project_id
    `

	sectionsJSON, _ := json.Marshal(p.Sections)
	now := time.Now()

	err := r.db.QueryRow(
		context.Background(),
		query,
		p.Title,
		p.Description,
		p.Introduction,
		p.OwnerID,
		now,
		string(sectionsJSON),
	).Scan(&p.ID)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return fmt.Errorf("project title already exists")
		}
		return fmt.Errorf("error creating project: %v", err)
	}

	p.CreatedAt = now
	p.UpdatedAt = now
	return nil
}

func (r *projectRepository) Update(p *models.Project) error {
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
		p.ID,
	)
	if err != nil {
		return fmt.Errorf("error updating project: %v", err)
	}

	p.UpdatedAt = now
	return nil
}

func (r *projectRepository) Delete(id string) error {
	result, err := r.db.Exec(context.Background(), "DELETE FROM projects WHERE project_id = $1", id)
	if err != nil {
		return fmt.Errorf("error deleting project: %v", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("project not found")
	}

	return nil
}
