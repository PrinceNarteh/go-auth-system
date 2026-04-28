package main

import (
	"auth-system/internal/utils"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
)

func initApp() *fiber.App {
	app := fiber.New(
		fiber.Config{
			ErrorHandler: utils.ErrorHandler,
		},
	)

	// Global middleware
	app.Use(logger.New())  // Log all requests
	app.Use(recover.New()) // Recover from panics
	app.Use(cors.New())    // Enable CORS for frontend apps

	return app
}
