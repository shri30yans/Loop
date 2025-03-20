 package repositories

import (
    "context"
    "time"
    "fmt"

    "Loop_backend/internal/models"
    "github.com/jackc/pgx/v4/pgxpool"
)

type postgresTagRepository struct {
    db *pgxpool.Pool
}

func NewPostgresTagRepository(db *pgxpool.Pool) TagRepository {
    return &postgresTagRepository{db: db}
}

func (r *postgresTagRepository) CreateTag(tag *models.Tag) error {
    now := time.Now()
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
        now,
        now,
    ).Scan(&tag.ID)
}

func (r *postgresTagRepository) GetTagByID(id string) (*models.Tag, error) {
    tag := &models.Tag{}
    query := `
        SELECT id, name, category, vector, usage_count, created_at, updated_at
        FROM tags
        WHERE id = $1
    `
    err := r.db.QueryRow(
        context.Background(),
        query,
        id,
    ).Scan(
        &tag.ID,
        &tag.Name,
        &tag.Category,
        &tag.Vector,
        &tag.UsageCount,
        &tag.CreatedAt,
        &tag.UpdatedAt,
    )
    if err != nil {
        return nil, fmt.Errorf("error getting tag by ID: %v", err)
    }
    return tag, nil
}

func (r *postgresTagRepository) GetTagByName(name string) (*models.Tag, error) {
    tag := &models.Tag{}
    query := `
        SELECT id, name, category, vector, usage_count, created_at, updated_at
        FROM tags
        WHERE name = $1
    `
    err := r.db.QueryRow(
        context.Background(),
        query,
        name,
    ).Scan(
        &tag.ID,
        &tag.Name,
        &tag.Category,
        &tag.Vector,
        &tag.UsageCount,
        &tag.CreatedAt,
        &tag.UpdatedAt,
    )
    if err != nil {
        return nil, fmt.Errorf("error getting tag by name: %v", err)
    }
    return tag, nil
}

func (r *postgresTagRepository) UpdateTag(tag *models.Tag) error {
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
    result, err := r.db.Exec(
        context.Background(),
        query,
        tag.Name,
        tag.Category,
        tag.Vector,
        tag.UsageCount,
        tag.UpdatedAt,
        tag.ID,
    )
    if err != nil {
        return fmt.Errorf("error updating tag: %v", err)
    }
    if result.RowsAffected() == 0 {
        return fmt.Errorf("no tag found with ID: %s", tag.ID)
    }
    return nil
}

func (r *postgresTagRepository) DeleteTag(id int) error {
    query := `DELETE FROM tags WHERE id = $1`
    result, err := r.db.Exec(context.Background(), query, id)
    if err != nil {
        return fmt.Errorf("error deleting tag: %v", err)
    }
    if result.RowsAffected() == 0 {
        return fmt.Errorf("no tag found with ID: %d", id)
    }
    return nil
}

func (r *postgresTagRepository) StoreTagVector(tagID int, vector []float64) error {
    query := `
        UPDATE tags
        SET vector = $1,
            updated_at = $2
        WHERE id = $3
    `
    result, err := r.db.Exec(
        context.Background(),
        query,
        vector,
        time.Now(),
        tagID,
    )
    if err != nil {
        return fmt.Errorf("error storing tag vector: %v", err)
    }
    if result.RowsAffected() == 0 {
        return fmt.Errorf("no tag found with ID: %d", tagID)
    }
    return nil
}

func (r *postgresTagRepository) FindSimilarTags(vector []float64, threshold float64, limit int) ([]*models.Tag, error) {
    query := `
        SELECT id, name, category, vector, usage_count, created_at, updated_at,
               1 - (vector <=> $1) as similarity
        FROM tags
        WHERE 1 - (vector <=> $1) > $2
        ORDER BY similarity DESC
        LIMIT $3
    `
    rows, err := r.db.Query(
        context.Background(),
        query,
        vector,
        threshold,
        limit,
    )
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
        tag.Confidence = similarity
        tags = append(tags, tag)
    }
    return tags, nil
}
