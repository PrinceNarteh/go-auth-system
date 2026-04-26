// Package handlers
package handlers

import "auth-system/internal/middleware"

type Handlers struct {
	Auth       AuthHandler
	Middleware middleware.Middleware
}
