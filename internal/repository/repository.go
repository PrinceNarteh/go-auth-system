// Package repository
package repository

import (
	"auth-system/internal/database"
)

type Repository struct {
	User UserRepository
}

func NewRepository(client *database.MongoDBClient) *Repository {
	return &Repository{
		User: &userRepository{collection: client.Database.Collection("users")},
	}
}
