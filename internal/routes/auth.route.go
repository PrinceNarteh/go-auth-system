package routes

import (
	"auth-system/internal/handlers"

	"github.com/gofiber/fiber/v3"
)

func authRoutes(r fiber.Router, h *handlers.Handlers) {
	r.Group("/auth")

	r.Post("/login", h.Auth.Login)
	r.Post("/register", h.Auth.Register)
}
