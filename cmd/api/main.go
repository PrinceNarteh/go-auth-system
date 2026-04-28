package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"auth-system/internal/config"
	"auth-system/internal/database"
	"auth-system/internal/handlers"
	"auth-system/internal/middleware"
	"auth-system/internal/repository"
	"auth-system/internal/routes"
	"auth-system/internal/services"
)

func main() {
	// connect MongoDB
	mongoClient, err := database.Connect()
	if err != nil {
		panic(err)
	}

	repo := repository.NewRepository(mongoClient)
	svc := services.NewService(repo)
	mw := middleware.NewMiddleware(svc)
	handlers := handlers.NewHandlers(svc, mw)

	app := initApp()

	routes.InitRoutes(app, handlers)

	// Listen from a different goroutine
	go func() {
		if err := app.Listen(config.Env.App.Port); err != nil {
			log.Panic(err)
		}
	}()

	quit := make(chan os.Signal, 1)                    // Create channel to signify a signal being sent
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM) // When an interrupt or termination signal is sent, notify the channel

	<-quit // This blocks the main thread until an interrupt is received
	fmt.Println("Gracefully shutting down...")
	err = app.Shutdown()
	if err != nil {
		log.Panic(err)
	}
}
