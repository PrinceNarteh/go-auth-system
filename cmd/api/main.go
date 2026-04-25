package main

import (
	"auth-system/internal/database"
)

func main() {
	// connect to database
	_, err := database.Connect()
	if err != nil {
		panic(err)
	}
}
