// Package services
package services

import "auth-system/internal/repository"

type Services struct {
	Auth AuthService
}

func NewService(repo *repository.Repository) *Services {
	return &Services{
		Auth: &authService{repo: repo},
	}
}
