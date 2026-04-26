// Package services
package services

import "auth-system/internal/repository"

type Services struct {
	Auth AuthService
	User UserService
}

func NewService(repo *repository.Repository) *Services {
	return &Services{
		Auth: &authService{repo: repo},
		User: &userService{repo: repo},
	}
}
