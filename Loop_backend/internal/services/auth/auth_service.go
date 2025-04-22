package auth

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"Loop_backend/internal/models"
	"Loop_backend/internal/repositories/auth"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	CreateSession(userID string) (*models.Session, error)
	ValidateToken(token string) (*models.Claims, error)
	AuthenticateUser(email string, password string) (string, error)
	RegisterUserPassword(userID, password string) error
}

type authService struct {
	secret   string
	repo     repositories.AuthRepository
	duration time.Duration
}

// NewAuthService creates a new authentication service
func NewAuthService(secret string, repo repositories.AuthRepository) AuthService {
	return &authService{
		secret:   secret,
		repo:     repo,
		duration: 24 * time.Hour * 100,
	}
}

func (s *authService) RegisterUserPassword(userID, password string) error {
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Call repository to store user and password
	err = s.repo.InsertUserPassword(userID, string(hashedPassword))
	if err != nil {
		return fmt.Errorf("failed to register user: %w", err)
	}

	return nil
}

func (s *authService) AuthenticateUser(email string, password string) (string, error) {

	user, err := s.repo.GetAuthenticatedUser(email)
	if err != nil {
		return "", fmt.Errorf("user not found")
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
	if err != nil {
		return "", fmt.Errorf("invalid password")
	}

	return user.UserID, nil
}

func (s *authService) CreateSession(userID string) (*models.Session, error) {
	// Create JWT claims
	claims := models.Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(s.secret))
	if err != nil {
		return nil, fmt.Errorf("error signing token: %v", err)
	}

	// Create session
	session := models.NewSession(userID, signedToken, s.duration)

	return session, nil
}

func (s *authService) ValidateToken(authHeader string) (*models.Claims, error) {

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		return nil, models.ErrMalformedToken
	}

	token, err := jwt.ParseWithClaims(tokenString, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.secret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, models.ErrExpiredToken
		}
		fmt.Println(err)
		return nil, models.ErrInvalidToken
	}

	claims, ok := token.Claims.(*models.Claims)
	if !ok || !token.Valid {
		return nil, models.ErrInvalidToken
	}

	if err := s.repo.CheckIfUserIdExists(claims.UserID); err != nil {
		fmt.Println(err)
		// User not found
		return nil, models.ErrInvalidToken
	}

	return claims, nil
}
