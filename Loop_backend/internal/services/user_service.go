package services

import (
    "Loop_backend/internal/models"
    "Loop_backend/internal/repositories"
)

type UserService interface {
    GetUser(user_id string) (*models.User, error)
    CreateUser(email, username string) (*models.User, error)
    UpdateUser(user_id string, email, username string) (*models.User, error)
    DeleteUser(user_id string) error
}

type userService struct {
    repo repositories.UserRepository
}

// NewUserService creates a new user service
func NewUserService(repo repositories.UserRepository) UserService {
    return &userService{repo: repo}
}

func (s *userService) GetUser(user_id string) (*models.User, error) {
    if user_id == "" {
        return nil, models.ErrInvalidID
    }
    return s.repo.GetUser(user_id)
}

func (s *userService) GetUserProjects(user_id string) (*models.User, error) {
    return s.repo.GetUser(user_id)
}

func (s *userService) CreateUser(email, username string) (*models.User, error) {

    // Create user instance
    newUser, err := models.NewUser(email, username, "", "")
    if err != nil {
        return nil, err
    }

    // Save to repository
    if err := s.repo.Create(newUser); err != nil {
        return nil, err
    }

    return newUser, nil
}

func (s *userService) UpdateUser(user_id string, email, username string) (*models.User, error) {
    // Get existing user
    existingUser, err := s.repo.GetUser(user_id)
    if err != nil {
        return nil, err
    }

    // Update fields
    if email != "" {
        existingUser.Email = email
    }
    if username != "" {
        existingUser.Username = username
    }

    // Save changes
    if err := s.repo.Update(existingUser); err != nil {
        return nil, err
    }

    return existingUser, nil
}

func (s *userService) DeleteUser(user_id string) error {
    if user_id == "" {
        return models.ErrInvalidID
    }
    return s.repo.Delete(user_id)
}

