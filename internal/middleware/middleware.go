package middleware

import "auth-system/internal/services"

type Middleware struct {
	Auth AuthMiddleware
}

func NewMiddleware(svc *services.Services) *Middleware {
	return &Middleware{
		Auth: AuthMiddleware{svc: svc},
	}
}
