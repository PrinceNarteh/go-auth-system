package services

import (
	"context"

	"auth-system/internal/models"
	"auth-system/internal/repository"
)

type UserService interface {
	GetUserByID(ctx context.Context, id string) (*models.User, error)
}

type userService struct {
	repo *repository.Repository
}

// GetUserByID retrieves a user by their ID
func (s *userService) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	return s.repo.User.FindByID(ctx, id)
}
