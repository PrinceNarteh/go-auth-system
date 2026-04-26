package repository

import (
	"context"
	"errors"
	"time"

	"auth-system/internal/models"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var (
	ErrEmailExists  = errors.New("email already exists")
	ErrUserNotFound = errors.New("user not found")
	ErrInvalidID    = errors.New("invalid id format")
)

type UserRepository interface {
	Create(context.Context, *models.User) error
	FindByEmail(context.Context, string) (*models.User, error)
	FindByID(context.Context, string) (*models.User, error)
	UpdateLastLogin(context.Context, string) error
}

type userRepository struct {
	collection *mongo.Collection
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	// Set timestamps and default values
	user.PrepareUser()

	res, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return ErrEmailExists
		}
		return err
	}
	user.ID = res.InsertedID.(bson.ObjectID)
	return nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	user := new(models.User)
	if err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

// FindByID finds a user by their MongoDB ObjectID
func (r *userRepository) FindByID(ctx context.Context, id string) (*models.User, error) {
	// Convert string ID to MongoDB ObjectID
	userID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, ErrInvalidID
	}
	user := new(models.User)
	if err := r.collection.FindOne(ctx, bson.M{"_id": userID}).Decode(user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

// UpdateLastLogin updates the user's last login timestamp
func (r *userRepository) UpdateLastLogin(ctx context.Context, id string) error {
	userID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return ErrInvalidID
	}
	filter := bson.M{"_id": userID}
	update := bson.M{
		"$set": bson.M{
			"lastLogin": time.Now(),
			"updatedAt": time.Now(),
		},
	}
	_, err = r.collection.UpdateOne(ctx, filter, update)
	return err
}
