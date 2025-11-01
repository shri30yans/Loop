package project

import (
    "context"
    "encoding/json"
    "fmt"
    "Loop_backend/internal/models"
    "github.com/google/uuid"
    "github.com/jackc/pgx/v4/pgxpool"
)

// ProjectRepository defines the interface for project data operations
type ProjectRepository interface {
    CreateProject(project *models.Project) error
    GetProject(projectID string) (*models.Project, error)
    UpdateProject(project *models.Project) error
    DeleteProject(projectID string) error
    SearchProjects(keyword string) ([]*models.ProjectInfo, error)
    GetProjectsByOwnerID(ownerID uuid.UUID) ([]*models.Project, error)
}

// PostgresProjectRepository implements ProjectRepository using PostgreSQL
type PostgresProjectRepository struct {
    pool *pgxpool.Pool
}

// NewProjectRepository creates a new project repository
func NewProjectRepository(pool *pgxpool.Pool) ProjectRepository {
    return &PostgresProjectRepository{pool: pool}
}

// CreateProject creates a new project
func (r *PostgresProjectRepository) CreateProject(project *models.Project) error {
    ctx := context.Background()
    
    // Convert sections to JSON
    sectionsJSON, err := json.Marshal(project.Sections)
    if err != nil {
        return fmt.Errorf("failed to marshal project sections: %w", err)
    }

    query := `
        INSERT INTO projects (
            project_id, 
            title, 
            description, 
            owner_id, 
            created_at, 
            updated_at,
            project_sections,
            status,
            introduction
        )
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
    
    _, err = r.pool.Exec(ctx, query,
        project.ProjectID,
        project.Title,
        project.Description,
        project.OwnerID,
        project.CreatedAt,
        project.UpdatedAt,
        sectionsJSON,
        project.Status,
        project.Introduction,
    )
    
    if err != nil {
        return fmt.Errorf("failed to create project: %w", err)
    }
    
    return nil
}

// GetProject retrieves a project by ID
func (r *PostgresProjectRepository) GetProject(projectID string) (*models.Project, error) {
    ctx := context.Background()
    query := `
        SELECT 
            project_id, 
            title, 
            description, 
            owner_id, 
            created_at, 
            updated_at,
            project_sections,
            status,
            introduction
        FROM projects
        WHERE project_id = $1`
    
    project := &models.Project{}
    var sectionsJSON []byte
    
    err := r.pool.QueryRow(ctx, query, projectID).Scan(
        &project.ProjectID,
        &project.Title,
        &project.Description,
        &project.OwnerID,
        &project.CreatedAt,
        &project.UpdatedAt,
        &sectionsJSON,
        &project.Status,
        &project.Introduction,
    )
    
    if err != nil {
        return nil, fmt.Errorf("failed to get project: %w", err)
    }
    
    // Unmarshal sections JSON
    if err := json.Unmarshal(sectionsJSON, &project.Sections); err != nil {
        return nil, fmt.Errorf("failed to unmarshal project sections: %w", err)
    }
    
    return project, nil
}

// UpdateProject updates an existing project
func (r *PostgresProjectRepository) UpdateProject(project *models.Project) error {
    ctx := context.Background()
    
    // Convert sections to JSON
    sectionsJSON, err := json.Marshal(project.Sections)
    if err != nil {
        return fmt.Errorf("failed to marshal project sections: %w", err)
    }

    query := `
        UPDATE projects 
        SET title = $1, 
            description = $2, 
            updated_at = $3,
            project_sections = $4,
            status = $5,
            introduction = $6
        WHERE project_id = $7`
    
    commandTag, err := r.pool.Exec(ctx, query,
        project.Title,
        project.Description,
        project.UpdatedAt,
        sectionsJSON,
        project.Status,
        project.Introduction,
        project.ProjectID,
    )
    
    if err != nil {
        return fmt.Errorf("failed to update project: %w", err)
    }
    
    if commandTag.RowsAffected() == 0 {
        return fmt.Errorf("no project found with ID %s", project.ProjectID)
    }
    
    return nil
}

// DeleteProject deletes a project by ID
func (r *PostgresProjectRepository) DeleteProject(projectID string) error {
    ctx := context.Background()
    query := `DELETE FROM projects WHERE project_id = $1`
    
    commandTag, err := r.pool.Exec(ctx, query, projectID)
    if err != nil {
        return fmt.Errorf("failed to delete project: %w", err)
    }
    
    if commandTag.RowsAffected() == 0 {
        return fmt.Errorf("no project found with ID %s", projectID)
    }
    
    return nil
}

// SearchProjects searches for projects by keyword
func (r *PostgresProjectRepository) SearchProjects(keyword string) ([]*models.ProjectInfo, error) {
    ctx := context.Background()
    query := `
        SELECT 
            project_id, 
            title, 
            description, 
            owner_id, 
            created_at,
            status,
            introduction
        FROM projects
        WHERE title ILIKE $1 OR description ILIKE $1
        ORDER BY created_at DESC`
    
    // If no keyword provided, return all projects
    searchPattern := "%"
    if keyword != "" {
        searchPattern = "%" + keyword + "%"
    }
    
    rows, err := r.pool.Query(ctx, query, searchPattern)
    if err != nil {
        return nil, fmt.Errorf("failed to search projects: %w", err)
    }
    defer rows.Close()
    
    var projects []*models.ProjectInfo
    for rows.Next() {
        project := &models.ProjectInfo{}
        err := rows.Scan(
            &project.ProjectID,
            &project.Title,
            &project.Description,
            &project.OwnerID,
            &project.CreatedAt,
            &project.Status,
            &project.Introduction,
        )
        if err != nil {
            return nil, fmt.Errorf("failed to scan project row: %w", err)
        }
        project.Tags = []string{};
        projects = append(projects, project)
    }
    
    if err = rows.Err(); err != nil {
        return nil, fmt.Errorf("error iterating project rows: %w", err)
    }
    
    return projects, nil
}

// GetProjectsByOwnerID retrieves all projects for a specific owner
func (r *PostgresProjectRepository) GetProjectsByOwnerID(ownerID uuid.UUID) ([]*models.Project, error) {
    ctx := context.Background()
    query := `
        SELECT 
            project_id, 
            title, 
            description, 
            owner_id, 
            created_at, 
            updated_at,
            project_sections,
            status,
            introduction
        FROM projects
        WHERE owner_id = $1
        ORDER BY created_at DESC`
    
    rows, err := r.pool.Query(ctx, query, ownerID)
    if err != nil {
        return nil, fmt.Errorf("failed to get projects by owner ID: %w", err)
    }
    defer rows.Close()
    
    var projects []*models.Project
    for rows.Next() {
        project := &models.Project{}
        var sectionsJSON []byte
        
        err := rows.Scan(
            &project.ProjectID,
            &project.Title,
            &project.Description,
            &project.OwnerID,
            &project.CreatedAt,
            &project.UpdatedAt,
            &sectionsJSON,
            &project.Status,
            &project.Introduction,
        )
        if err != nil {
            return nil, fmt.Errorf("failed to scan project row: %w", err)
        }
        
        // Unmarshal sections JSON
        if err := json.Unmarshal(sectionsJSON, &project.Sections); err != nil {
            return nil, fmt.Errorf("failed to unmarshal project sections: %w", err)
        }
        
        projects = append(projects, project)
    }
    
    if err = rows.Err(); err != nil {
        return nil, fmt.Errorf("error iterating project rows: %w", err)
    }
    
    return projects, nil
}
