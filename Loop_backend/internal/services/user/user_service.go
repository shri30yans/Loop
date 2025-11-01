package user

import (
	"Loop_backend/internal/models"
	userRepo "Loop_backend/internal/repositories/user"
)

type UserService interface {
	GetUser(user_id string) (*models.UserInfo, error)
	CreateUser(email, username string) (*models.User, error)
	UpdateUser(user_id string, email, username string) (*models.User, error)
	DeleteUser(user_id string) error
}

type userService struct {
repo userRepo.UserRepository
}

func NewUserService(repo userRepo.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) GetUser(user_id string) (*models.UserInfo, error) {
	if user_id == "" {
		return nil, models.ErrInvalidID
	}
	return s.repo.GetUser(user_id)
}

func (s *userService) CreateUser(email, username string) (*models.User, error) {

	newUser, err := models.NewUser(email, username, "", "")
	if err != nil {
		return nil, err
	}

	if err := s.repo.Create(newUser); err != nil {
		return nil, err
	}

	return newUser, nil
}

func (s *userService) UpdateUser(user_id string, email, username string) (*models.User, error) {
	userInfo, err := s.repo.GetUser(user_id)
	if err != nil {
		return nil, err
	}

	if email != "" {
		userInfo.Email = email
	}
	if username != "" {
		userInfo.Username = username
	}

	if err := s.repo.Update(&userInfo.User); err != nil {
		return nil, err
	}

	return &userInfo.User, nil
}

func (s *userService) DeleteUser(user_id string) error {
	if user_id == "" {
		return models.ErrInvalidID
	}
	return s.repo.Delete(user_id)
}
