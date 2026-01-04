package services

import (
	"Golang-backend/internal/model"
	"Golang-backend/internal/repository"
	"errors"
	"fmt"
)

type UserService struct {
	userRepo *repository.UserRepository
}

var ErrUserNotFound = errors.New("user not found")

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) CreateUser(user model.User) (int, error) {
	id, err := s.userRepo.Create(user)
	if err != nil {
		if err.Error() == "pq: duplicate key value violates unique constraint \"user_email_key\"" {
			return 0, fmt.Errorf("Email already exists")
		}
		return 0, fmt.Errorf("%s", err)
	}
	return id, nil
}

func (s *UserService) GetUser(id int64) (*model.User, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("Failed to get user by ID %d: %w", id, err)
	}
	if user == nil {
		return nil, fmt.Errorf("%w: id=%d", ErrUserNotFound, id)
	}
	return user, nil
}

func (s *UserService) UpdateUser(id int64, data model.UserUpdate) (*model.User, error) {
	user, err := s.GetUser(id)
	if err != nil {
		return nil, fmt.Errorf("Failed to update user: id: %d error: %w", id, err)
	}
	if data.Name != "" {
		user.Name = data.Name
	}
	if data.Email != "" {
		user.Email = data.Email
	}
	new_data, err := s.userRepo.UpdateByID(id, *user)
	if err != nil {
		return nil, fmt.Errorf("Failed to update user in DB: %w", err)
	}
	return new_data, nil
}
