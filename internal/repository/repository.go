// Package repository
package repository

import "go.mongodb.org/mongo-driver/v2/mongo"

type Repository struct {
	User UserRepository
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		User: &userRepository{collection: db.Collection("user")},
	}
}
