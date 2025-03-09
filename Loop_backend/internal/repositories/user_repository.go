package repositories

import (
    "context"
    "errors"
    "fmt"
    "time"

    "github.com/jackc/pgconn"
    "github.com/jackc/pgx/v4/pgxpool"
    "Loop_backend/internal/models"
)

type UserRepository interface {
    FindByID(user_id string) (*models.User, error)
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

func (r *userRepository) FindByID(id string) (*models.User, error) {
    query := `
    SELECT id, email, username, bio, location, created_at, updated_at
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
        &u.UpdatedAt,
    )

    if err != nil {
        if err.Error() == "no rows in result set" {
            return nil, fmt.Errorf("user not found: %v", err)
        }
        return nil, fmt.Errorf("error finding user: %v", err)
    }

    return &u, nil
}

func (r *userRepository) Create(u *models.User) error {
    query := `
    INSERT INTO users (email, username, bio, location, created_at)
    VALUES ($1, $2, $3, $4, $5)
    RETURNING id
    `

    err := r.db.QueryRow(
        context.Background(),
        query,
        u.Email,
        u.Username,
        u.Bio,
        u.Location,
        time.Now(),
    ).Scan(&u.ID)

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
    SET email = $1, username = $2, bio = $3, location = $4, updated_at = $5
    WHERE id = $6
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