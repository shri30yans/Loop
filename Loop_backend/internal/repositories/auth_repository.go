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
}


type authRepository struct {
	db *pgxpool.Pool

}

// NewAuthRepository creates a new PostgreSQL auth repository
func NewAuthRepository(db *pgxpool.Pool) *authRepository {
	return &authRepository{db: db}
}


func (r *authRepository) GetAuthenticatedUser(email string) (*models.AuthenticatedUser, error) {
    query := `
    SELECT users.user_id, hashed_password
    FROM passwords
    JOIN users
    ON passwords.user_id = users.user_id
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

