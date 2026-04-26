package services

import (
	"context"

	"auth-system/internal/dtos"
	"auth-system/internal/models"
	"auth-system/internal/repository"
)

var _ AuthService = (*authService)(nil)

type AuthService interface {
	Register(context.Context, *dtos.RegisterDTO) (*models.User, error)
	Login(context.Context, *dtos.LoginDTO) (*models.LoginResponse, error)
}

type authService struct {
	repo *repository.Repository
}

func (s *authService) Login(ctx context.Context, data *dtos.LoginDTO) (*models.LoginResponse, error) {
	return nil, nil
}

func (s *authService) Register(ctx context.Context, data *dtos.RegisterDTO) (*models.User, error) {
	return nil, nil
}
