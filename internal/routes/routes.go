// Package routes
package routes

import (
	"time"

	"auth-system/internal/handlers"

	"github.com/gofiber/fiber/v3"
)

func InitRoutes(app fiber.Router, h *handlers.Handlers) {
	// Health check endpoint (useful for container orchestration)
	app.Get("/health", func(c fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":    "ok",
			"timestamp": time.Now().Unix(),
		})
	})

	// /api
	api := app.Group("/api")

	// auth routes
	authRoutes(api, h)
}
