package repositories

import (
	"context"
	"fmt"
	"time"

	"Loop_backend/internal/models"
	"github.com/jackc/pgx/v4/pgxpool"
)

type tagRepository struct {
	db *pgxpool.Pool
}

func NewTagRepository(db *pgxpool.Pool) TagRepository {
	return &tagRepository{db: db}
}

func (r *tagRepository) CreateTag(tag *models.Tag) error {
	query := `
		INSERT INTO tags (name, category, vector, usage_count, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	return r.db.QueryRow(
		context.Background(),
		query,
		tag.Name,
		tag.Category,
		tag.Vector,
		tag.UsageCount,
		tag.CreatedAt,
		tag.UpdatedAt,
	).Scan(&tag.ID)
}

func (r *tagRepository) GetTagByID(id string) (*models.Tag, error) {
	query := `
		SELECT id, name, category, vector, usage_count, created_at, updated_at
		FROM tags
		WHERE id = $1
	`
	tag := &models.Tag{}
	err := r.db.QueryRow(context.Background(), query, id).Scan(
		&tag.ID,
		&tag.Name,
		&tag.Category,
		&tag.Vector,
		&tag.UsageCount,
		&tag.CreatedAt,
		&tag.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("error getting tag: %v", err)
	}
	return tag, nil
}

func (r *tagRepository) GetTagByName(name string) (*models.Tag, error) {
	query := `
		SELECT id, name, category, vector, usage_count, created_at, updated_at
		FROM tags
		WHERE name = $1
	`
	tag := &models.Tag{}
	err := r.db.QueryRow(context.Background(), query, name).Scan(
		&tag.ID,
		&tag.Name,
		&tag.Category,
		&tag.Vector,
		&tag.UsageCount,
		&tag.CreatedAt,
		&tag.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("error getting tag: %v", err)
	}
	return tag, nil
}

func (r *tagRepository) UpdateTag(tag *models.Tag) error {
	query := `
		UPDATE tags
		SET name = $1,
			category = $2,
			vector = $3,
			usage_count = $4,
			updated_at = $5
		WHERE id = $6
	`
	tag.UpdatedAt = time.Now()
	_, err := r.db.Exec(
		context.Background(),
		query,
		tag.Name,
		tag.Category,
		tag.Vector,
		tag.UsageCount,
		tag.UpdatedAt,
		tag.ID,
	)
	return err
}

func (r *tagRepository) DeleteTag(id int) error {
	_, err := r.db.Exec(context.Background(), "DELETE FROM tags WHERE id = $1", id)
	return err
}

func (r *tagRepository) StoreTagVector(tagID int, vector []float64) error {
	_, err := r.db.Exec(
		context.Background(),
		"UPDATE tags SET vector = $1, updated_at = $2 WHERE id = $3",
		vector,
		time.Now(),
		tagID,
	)
	return err
}

func (r *tagRepository) FindSimilarTags(vector []float64, threshold float64, limit int) ([]*models.Tag, error) {
	query := `
		SELECT id, name, category, vector, usage_count, created_at, updated_at,
			   1 - (vector <=> $1) as similarity
		FROM tags
		WHERE 1 - (vector <=> $1) > $2
		ORDER BY similarity DESC
		LIMIT $3
	`
	rows, err := r.db.Query(context.Background(), query, vector, threshold, limit)
	if err != nil {
		return nil, fmt.Errorf("error finding similar tags: %v", err)
	}
	defer rows.Close()

	var tags []*models.Tag
	for rows.Next() {
		tag := &models.Tag{}
		var similarity float64
		err := rows.Scan(
			&tag.ID,
			&tag.Name,
			&tag.Category,
			&tag.Vector,
			&tag.UsageCount,
			&tag.CreatedAt,
			&tag.UpdatedAt,
			&similarity,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning tag row: %v", err)
		}
		tag.SetConfidence(similarity)
		tags = append(tags, tag)
	}
	return tags, nil
}

func (r *tagRepository) CreateTagRelationship(tag1ID, tag2ID int, strength float64) error {
	query := `
		INSERT INTO tag_relationships (tag1_id, tag2_id, strength, co_occurrences, last_updated)
		VALUES ($1, $2, $3, 1, $4)
		ON CONFLICT (tag1_id, tag2_id) DO UPDATE
		SET strength = tag_relationships.strength + $3,
			co_occurrences = tag_relationships.co_occurrences + 1,
			last_updated = $4
	`
	_, err := r.db.Exec(
		context.Background(),
		query,
		tag1ID,
		tag2ID,
		strength,
		time.Now(),
	)
	return err
}

func (r *tagRepository) UpdateTagRelationship(tag1ID, tag2ID int, strength float64) error {
	_, err := r.db.Exec(
		context.Background(),
		`UPDATE tag_relationships 
		 SET strength = $3, last_updated = $4
		 WHERE tag1_id = $1 AND tag2_id = $2`,
		tag1ID,
		tag2ID,
		strength,
		time.Now(),
	)
	return err
}

func (r *tagRepository) GetRelatedTags(tagID int, minStrength float64) ([]*models.Tag, error) {
	query := `
		SELECT t.id, t.name, t.category, t.vector, t.usage_count, t.created_at, t.updated_at,
			   tr.strength as confidence
		FROM tags t
		JOIN tag_relationships tr ON (tr.tag2_id = t.id)
		WHERE tr.tag1_id = $1 AND tr.strength >= $2
		ORDER BY tr.strength DESC
	`
	rows, err := r.db.Query(context.Background(), query, tagID, minStrength)
	if err != nil {
		return nil, fmt.Errorf("error getting related tags: %v", err)
	}
	defer rows.Close()

	var tags []*models.Tag
	for rows.Next() {
		tag := &models.Tag{}
		err := rows.Scan(
			&tag.ID,
			&tag.Name,
			&tag.Category,
			&tag.Vector,
			&tag.UsageCount,
			&tag.CreatedAt,
			&tag.UpdatedAt,
			&tag.Confidence,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning tag row: %v", err)
		}
		tags = append(tags, tag)
	}
	return tags, nil
}

func (r *tagRepository) AssignTagToProject(projectID string, tagID int, confidence float64) error {
	query := `
		INSERT INTO project_tags (project_id, tag_description, confidence, created_at)
		VALUES ($1, (SELECT name FROM tags WHERE id = $2), $3, $4)
		ON CONFLICT (project_id, tag_description)
		DO UPDATE SET confidence = $3
	`
	_, err := r.db.Exec(
		context.Background(),
		query,
		projectID,
		tagID,
		confidence,
		time.Now(),
	)
	return err
}

func (r *tagRepository) GetProjectTags(projectID string) ([]*models.Tag, error) {
	query := `
		SELECT t.id, t.name, t.category, t.vector, t.usage_count, t.created_at, t.updated_at,
			   pt.confidence
		FROM tags t
		JOIN project_tags pt ON pt.tag_description = t.name
		WHERE pt.project_id = $1
		ORDER BY pt.confidence DESC
	`
	rows, err := r.db.Query(context.Background(), query, projectID)
	if err != nil {
		return nil, fmt.Errorf("error getting project tags: %v", err)
	}
	defer rows.Close()

	var tags []*models.Tag
	for rows.Next() {
		tag := &models.Tag{}
		err := rows.Scan(
			&tag.ID,
			&tag.Name,
			&tag.Category,
			&tag.Vector,
			&tag.UsageCount,
			&tag.CreatedAt,
			&tag.UpdatedAt,
			&tag.Confidence,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning tag row: %v", err)
		}
		tags = append(tags, tag)
	}
	return tags, nil
}
