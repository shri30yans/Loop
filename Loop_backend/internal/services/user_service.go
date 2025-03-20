package services

import (
    "Loop_backend/internal/models"
    "Loop_backend/internal/repositories"
)

type UserService interface {
<<<<<<< HEAD
    GetUser(user_id string) (*models.UserInfo, error)
=======
    GetUser(user_id string) (*models.User, error)
>>>>>>> 4a2f436bed91636c5c2e3782993f5ab211ecfca7
    CreateUser(email, username string) (*models.User, error)
    UpdateUser(user_id string, email, username string) (*models.User, error)
    DeleteUser(user_id string) error
}

type userService struct {
    repo repositories.UserRepository
}

<<<<<<< HEAD
=======
// NewUserService creates a new user service
>>>>>>> 4a2f436bed91636c5c2e3782993f5ab211ecfca7
func NewUserService(repo repositories.UserRepository) UserService {
    return &userService{repo: repo}
}

<<<<<<< HEAD
func (s *userService) GetUser(user_id string) (*models.UserInfo, error) {
=======
func (s *userService) GetUser(user_id string) (*models.User, error) {
>>>>>>> 4a2f436bed91636c5c2e3782993f5ab211ecfca7
    if user_id == "" {
        return nil, models.ErrInvalidID
    }
    return s.repo.GetUser(user_id)
}

<<<<<<< HEAD

func (s *userService) CreateUser(email, username string) (*models.User, error) {

=======
func (s *userService) GetUserProjects(user_id string) (*models.User, error) {
    return s.repo.GetUser(user_id)
}

func (s *userService) CreateUser(email, username string) (*models.User, error) {

    // Create user instance
>>>>>>> 4a2f436bed91636c5c2e3782993f5ab211ecfca7
    newUser, err := models.NewUser(email, username, "", "")
    if err != nil {
        return nil, err
    }

<<<<<<< HEAD
=======
    // Save to repository
>>>>>>> 4a2f436bed91636c5c2e3782993f5ab211ecfca7
    if err := s.repo.Create(newUser); err != nil {
        return nil, err
    }

    return newUser, nil
}

func (s *userService) UpdateUser(user_id string, email, username string) (*models.User, error) {
<<<<<<< HEAD
    userInfo, err := s.repo.GetUser(user_id)
=======
    // Get existing user
    existingUser, err := s.repo.GetUser(user_id)
>>>>>>> 4a2f436bed91636c5c2e3782993f5ab211ecfca7
    if err != nil {
        return nil, err
    }

<<<<<<< HEAD
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
=======
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
>>>>>>> 4a2f436bed91636c5c2e3782993f5ab211ecfca7
}

func (s *userService) DeleteUser(user_id string) error {
    if user_id == "" {
        return models.ErrInvalidID
    }
    return s.repo.Delete(user_id)
}
<<<<<<< HEAD
=======

>>>>>>> 4a2f436bed91636c5c2e3782993f5ab211ecfca7
