package repositories

import (
	"context"
	"fmt"
	"errors"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
    "Loop_backend/internal/models"
)

type AuthRepository interface {
	GetAuthenticatedUser(email string) (*models.AuthenticatedUser, error)
    InsertUserPassword(userID string, hashedPassword string) error
    CheckIfUserIdExists(userID string) error
}


type authRepository struct {
	db *pgxpool.Pool

}

// NewAuthRepository creates a new PostgreSQL auth repository
func NewAuthRepository(db *pgxpool.Pool) *authRepository {
	return &authRepository{db: db}
}

func (r *authRepository) InsertUserPassword(userID string, hashedPassword string) error {
    query := `
    INSERT INTO passwords(user_id, hashed_password)
    VALUES($1, $2)
    `

    _, err := r.db.Exec(context.Background(), query, userID, hashedPassword)
    if err != nil {
        return fmt.Errorf("failed to insert user password: %w", err)
    }

    return nil
}

func (r *authRepository) CheckIfUserIdExists(id string) error {
    query := `
    SELECT users.id
    FROM users
    WHERE users.id = $1
    `
    var user_id string;
    err := r.db.QueryRow(context.Background(), query, id).Scan(
        &user_id,
    )

    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return fmt.Errorf("user not found: %w", err)
        }
        return fmt.Errorf("error finding user: %w", err)
    }

    return nil
}


func (r *authRepository) GetAuthenticatedUser(email string) (*models.AuthenticatedUser, error) {
    query := `
    SELECT users.id, hashed_password
    FROM passwords
    JOIN users
    ON passwords.user_id = users.id
    WHERE users.email = $1
    LIMIT 1
    `

    var user models.AuthenticatedUser
    err := r.db.QueryRow(context.Background(), query, email).Scan(
        &user.UserID,
        &user.HashedPassword,
    )

    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, fmt.Errorf("user not found: %w", err)
        }
        return nil, fmt.Errorf("error finding user: %w", err)
    }

    return &user, nil
}




