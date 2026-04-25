// Package database
package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"auth-system/internal/config"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const QueryTimeout = 10 * time.Second

type MongoDBClient struct {
	Client   *mongo.Client
	Database *mongo.Database
}

// Connect establishes connection to MongoDB
// This function handles connection pooling and error management
func Connect() (*MongoDBClient, error) {
	// Create context with timeout for connection
	// 10 seconds should be enough for connection establishment
	ctx, cancel := context.WithTimeout(context.Background(), QueryTimeout)
	defer cancel()

	clientOptions := options.Client().
		ApplyURI(config.Env.DB.URI).
		SetMinPoolSize(10).
		SetMaxPoolSize(100).
		SetTimeout(30 * time.Second).
		SetMaxConnIdleTime(30 * time.Second).
		SetAuth(options.Credential{
			Username: config.Env.DB.Username,
			Password: config.Env.DB.Password,
		})

	client, err := mongo.Connect(clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Ping the database to verify connection
	// This ensures we can actually communicate with the database
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	log.Println("Successfully connected to MongoDB")

	return &MongoDBClient{
		Client:   client,
		Database: client.Database(config.Env.DB.Name),
	}, nil
}

// Close gracefully closes the MongoDB connection
// Should be called when shutting down the application
func (m *MongoDBClient) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), QueryTimeout)
	defer cancel()

	if err := m.Client.Disconnect(ctx); err != nil {
		return fmt.Errorf("failed to disconnect MongoDB: %w", err)
	}

	log.Println("MongoDB connection closed.")
	return nil
}

// GetCollection returns a MongoDB collection
// This is a helper to avoid typing the full path every time
func (m *MongoDBClient) GetCollection(name string) *mongo.Collection {
	return m.Database.Collection(name)
}
