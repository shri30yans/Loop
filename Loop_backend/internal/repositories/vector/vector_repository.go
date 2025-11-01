package vector

import (
    "context"
    "fmt"
    "strconv"
    "strings"
    "github.com/google/uuid"
    "github.com/jackc/pgx/v5/pgxpool"
)

// VectorRepository defines the interface for vector operations
type VectorRepository interface {
    StoreProjectVector(projectID uuid.UUID, summary string, vector []float32) error
    GetProjectVector(projectID uuid.UUID) ([]float32, error)
    SearchSimilarProjects(vector []float32, limit int) ([]uuid.UUID, error)
}

// PostgresVectorRepository implements VectorRepository using PostgreSQL
type PostgresVectorRepository struct {
    pool *pgxpool.Pool
}

// InitDB initializes a connection to the database
func InitDB(connString string) (*pgxpool.Pool, error) {
    ctx := context.Background()
    return pgxpool.New(ctx, connString)
}

// NewVectorRepository creates a new vector repository
func NewVectorRepository(pool *pgxpool.Pool) VectorRepository {
    return &PostgresVectorRepository{
        pool: pool,
    }
}

// StoreProjectVector stores a project's vector representation and summary
func (r *PostgresVectorRepository) StoreProjectVector(projectID uuid.UUID, summary string, vector []float32) error {
    ctx := context.Background()
    // Convert []float32 to a string array format that pgvector expects
    vectorStr := "["
    for i, v := range vector {
        if i > 0 {
            vectorStr += ","
        }
        vectorStr += fmt.Sprintf("%f", v)
    }
    vectorStr += "]"

    query := `
        INSERT INTO project_vectors (project_id, summary, embedding)
        VALUES ($1, $2, $3::vector)
        ON CONFLICT (project_id)
        DO UPDATE SET summary = EXCLUDED.summary, embedding = EXCLUDED.embedding`

    _, err := r.pool.Exec(ctx, query, projectID, summary, vectorStr)
    return err
}

// GetProjectVector retrieves a project's vector representation
func (r *PostgresVectorRepository) GetProjectVector(projectID uuid.UUID) ([]float32, error) {
    ctx := context.Background()
    query := `
        SELECT embedding::text
        FROM project_vectors
        WHERE project_id = $1`

    var vectorStr string
    err := r.pool.QueryRow(ctx, query, projectID).Scan(&vectorStr)
    if err != nil {
        return nil, err
    }

    // Parse the vector string back into []float32
    // Remove brackets
    vectorStr = vectorStr[1 : len(vectorStr)-1]
    // Split by comma
    parts := strings.Split(vectorStr, ",")
    vector := make([]float32, len(parts))
    for i, part := range parts {
        val, err := strconv.ParseFloat(strings.TrimSpace(part), 32)
        if err != nil {
            return nil, fmt.Errorf("failed to parse vector component: %w", err)
        }
        vector[i] = float32(val)
    }

    return vector, nil
}

// SearchSimilarProjects finds similar projects based on vector similarity
func (r *PostgresVectorRepository) SearchSimilarProjects(vector []float32, limit int) ([]uuid.UUID, error) {
    ctx := context.Background()
    // Convert input vector to pgvector format
    vectorStr := "["
    for i, v := range vector {
        if i > 0 {
            vectorStr += ","
        }
        vectorStr += fmt.Sprintf("%f", v)
    }
    vectorStr += "]"

    query := `
        SELECT project_id
        FROM project_vectors
        ORDER BY embedding <-> $1::vector
        LIMIT $2`

    rows, err := r.pool.Query(ctx, query, vectorStr, limit)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var projectIDs []uuid.UUID
    for rows.Next() {
        var projectID uuid.UUID
        if err := rows.Scan(&projectID); err != nil {
            return nil, err
        }
        projectIDs = append(projectIDs, projectID)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

    return projectIDs, nil
}
