package repositories

import (
	"Loop_backend/internal/models"
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type TagRepository struct {
	db *pgxpool.Pool
}

func NewTagRepository(db *pgxpool.Pool) *TagRepository {
	return &TagRepository{db: db}
}

// CreateTag creates a new tag in the database
func (r *TagRepository) CreateTag(tag *models.Tag) error {
	tag.BeforeCreate()
	query := `
        INSERT INTO tags (
            id, project_id, name, type, description, 
            category, usage_count, embedding, created_at, updated_at
        ) VALUES (
            $1, $2, $3, $4, $5, 
            $6, $7, $8, $9, $10
        )
    `
	_, err := r.db.Exec(
		context.Background(),
		query,
		tag.ID,
		tag.ProjectID,
		tag.Name,
		tag.Type,
		tag.Description,
		tag.Category,
		tag.UsageCount,
		tag.Embedding,
		tag.CreatedAt,
		tag.UpdatedAt,
	)
	return err
}

// GetTagsByProjectID retrieves all tags for a project
func (r *TagRepository) GetTagsByProjectID(projectID uuid.UUID) ([]models.Tag, error) {
	query := `
        SELECT 
            id, project_id, name, type, description,
            category, usage_count, embedding, created_at, updated_at
        FROM tags 
        WHERE project_id = $1
    `
	rows, err := r.db.Query(context.Background(), query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []models.Tag
	for rows.Next() {
		var tag models.Tag
		err := rows.Scan(
			&tag.ID,
			&tag.ProjectID,
			&tag.Name,
			&tag.Type,
			&tag.Description,
			&tag.Category,
			&tag.UsageCount,
			&tag.Embedding,
			&tag.CreatedAt,
			&tag.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, rows.Err()
}
