// Package handlers
package handlers

import (
	"auth-system/internal/middleware"
	"auth-system/internal/services"
)

type Handlers struct {
	Auth       AuthHandler
	Middleware *middleware.Middleware
}

func NewHandlers(svc *services.Services, middleware *middleware.Middleware) *Handlers {
	return &Handlers{
		Auth:       &authHandler{svc: svc},
		Middleware: middleware,
	}
}
