// Package middleware
package middleware

import (
	"slices"
	"strings"

	"auth-system/internal/services"
	"auth-system/internal/utils"

	"github.com/gofiber/fiber/v3"
)

type AuthMiddleware struct {
	svc *services.Services
}

// AuthRequired protects routes that require authentication
// It extracts and validates the JWT token from the Authorization header
func (m *AuthMiddleware) AuthRequired() fiber.Handler {
	return func(c fiber.Ctx) error {
		// Get the Authorization header
		authHeader := c.Get("authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "missing authorization header",
			})
		}

		// Check if header has Bearer prefix
		// Format should be: "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid authorization header format. Use: Bearer <token>",
			})
		}

		token := parts[1]

		// Validate the token and get user
		user, err := utils.ValidateToken(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token: " + err.Error(),
			})
		}

		// Store user in context for later use in handlers
		c.Locals("user", user)
		c.Locals("userID", user.ID)
		c.Locals("userRole", user.Role)

		// Proceed to next middleware/handler
		return c.Next()
	}
}

// RolesMiddleware ensures the user has admin role
// Must be used after AuthMiddleware
func (m *AuthMiddleware) RolesMiddleware(roles ...string) fiber.Handler {
	return func(c fiber.Ctx) error {
		role, ok := c.Locals("userRole").(string)
		if !ok || !slices.Contains(roles, role) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Amdin access required",
			})
		}
		return c.Next()
	}
}

// OptionalAuthMiddleware tries to authenticate but doesn't require it
// Useful for routes that work with or without authentication
func (m *AuthMiddleware) OptionalAuthMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		authHeader := c.Get("authorization")
		if authHeader == "" {
			return c.Next()
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Next()
		}

		token := parts[1]
		user, err := m.svc.Auth.ValidateToken(token)
		if err == nil {
			c.Locals("user", user)
			c.Locals("userID", user.ID.Hex())
		}

		return c.Next()
	}
}

// RefreshToken generates a new token from an old one
// POST /api/auth/refresh
func (m *AuthMiddleware) RefreshToken(c fiber.Ctx) error {
	// Get token from Authorization header
	authHeader := c.Get("authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Mission token",
		})
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid token format",
		})
	}

	// Refresh token
	oldToken := parts[1]
	newToken, err := utils.RefreshToken(oldToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Failed to refresh token: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"token": newToken,
	})
}
