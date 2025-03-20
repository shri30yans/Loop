package repositories

import (
"context"
"errors"
"fmt"
"time"

"Loop_backend/internal/models"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
)

type UserRepository interface {
GetUser(user_id string) (*models.UserInfo, error)
Create(user *models.User) error
Update(user *models.User) error
Delete(user_id string) error
}

type userRepository struct {
	db *pgxpool.Pool
}

// NewUserRepository creates a new PostgreSQL user repository
func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetUser(id string) (*models.UserInfo, error) {
    // First fetch the user
    query := `
        SELECT id, email, username, bio, location, created_at
        FROM users
        WHERE id = $1
    `

    var u models.User
    err := r.db.QueryRow(context.Background(), query, id).Scan(
        &u.ID,
        &u.Email,
        &u.Username,
        &u.Bio,
        &u.Location,
        &u.CreatedAt,
    )

    if err != nil {
        if err.Error() == "no rows in result set" {
            return nil, fmt.Errorf("user not found: %v", err)
        }
        return nil, fmt.Errorf("error finding user: %v", err)
    }

    // Then fetch all projects owned by this user
    projectsQuery := `
        SELECT p.project_id, p.title, p.description, p.status, p.introduction, p.owner_id, 
               p.created_at, p.updated_at, p.project_sections::TEXT
        FROM projects p
        WHERE p.owner_id = $1
        ORDER BY p.created_at DESC
    `

    rows, err := r.db.Query(context.Background(), projectsQuery, id)
    if err != nil {
        return &models.UserInfo{User: u}, fmt.Errorf("error fetching user projects: %v", err)
    }
    defer rows.Close()

    var projectInfos []*models.ProjectInfo
    for rows.Next() {
        var p models.ProjectInfo
        var sectionsJSON string

        err := rows.Scan(
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
            return &models.UserInfo{User: u}, fmt.Errorf("error scanning project row: %v", err)
        }

        // Initialize empty tags array
        p.Tags = []string{}

        projectInfos = append(projectInfos, &p)
    }

    return &models.UserInfo{
        User:     u,
        Projects: projectInfos,
    }, nil
}

func (r *userRepository) Create(u *models.User) error {
	query := `
    INSERT INTO users (id, email, username, bio, location, created_at)
    VALUES ($1, $2, $3, $4, $5, $6)
    `
	_, err := r.db.Exec(
		context.Background(),
		query,
		u.ID,
		u.Email,
		u.Username,
		u.Bio,
		u.Location,
		time.Now(),
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return fmt.Errorf("email already exists")
		}
		return fmt.Errorf("error creating user: %v", err)
	}

	return nil
}

func (r *userRepository) Update(u *models.User) error {
	query := `
    UPDATE users
    SET email = $1, username = $2, bio = $3, location = $4
    WHERE id = $5
    `

	commandTag, err := r.db.Exec(
		context.Background(),
		query,
		u.Email,
		u.Username,
		u.Bio,
		u.Location,
		time.Now(),
		u.ID,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return fmt.Errorf("email already exists")
		}
		return fmt.Errorf("error updating user: %v", err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (r *userRepository) Delete(id string) error {
	query := `DELETE FROM users WHERE id = $1`

	commandTag, err := r.db.Exec(context.Background(), query, id)
	if err != nil {
		return fmt.Errorf("error deleting user: %v", err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}
