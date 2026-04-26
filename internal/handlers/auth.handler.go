package handlers

import (
	"errors"

	"auth-system/internal/dtos"
	"auth-system/internal/models"
	"auth-system/internal/repository"
	"auth-system/internal/services"

	"github.com/gofiber/fiber/v3"
)

var _ AuthHandler = (*authHandler)(nil)

type AuthHandler interface {
	Login(fiber.Ctx) error
	Register(fiber.Ctx) error
	GetProfile(fiber.Ctx) error
}

type authHandler struct {
	svc *services.Services
}

// Register handles user registration
// POST /api/auth/register
func (h *authHandler) Register(c fiber.Ctx) error {
	registerDTO := new(dtos.RegisterDTO)
	if err := c.Bind().Body(registerDTO); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	user, err := h.svc.Auth.Register(c.Context(), registerDTO)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrEmailExists):
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "Email already registered",
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to register user: " + err.Error(),
			})
		}
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User registered successfully",
		"user":    user,
	})
}

// Login handles user login
// POST /api/auth/login
func (h *authHandler) Login(c fiber.Ctx) error {
	loginDTO := new(dtos.LoginDTO)
	if err := c.Bind().Body(loginDTO); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	user, err := h.svc.Auth.Login(c.Context(), loginDTO)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Login successful",
		"data":    user,
	})
}

// GetProfile returns the authenticated user's profile
// GET /api/auth/profile
func (h *authHandler) GetProfile(c fiber.Ctx) error {
	// Get user from context (set by auth middleware)
	userRaw := c.Locals("user")
	if userRaw == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not found in context",
		})
	}

	user, ok := userRaw.(*models.User)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Invalid user data in context",
		})
	}

	// Remove sensitive data
	user.Password = ""

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"user": user,
	})
}
