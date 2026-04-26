// Package routes
package routes

import (
	"auth-system/internal/handlers"

	"github.com/gofiber/fiber/v3"
)

func InitRoutes(router fiber.Router, h *handlers.Handlers) {
	r := router.Group("/api")
	authRoutes(r, h)
}
