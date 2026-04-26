package services

import (
	"context"
	"errors"
	"fmt"

	"auth-system/internal/dtos"
	"auth-system/internal/models"
	"auth-system/internal/repository"
	"auth-system/internal/utils"

	"go.mongodb.org/mongo-driver/v2/bson"
)

var _ AuthService = (*authService)(nil)

type AuthService interface {
	Register(context.Context, *dtos.RegisterDTO) (*models.User, error)
	Login(context.Context, *dtos.LoginDTO) (*models.LoginResponse, error)
}

type authService struct {
	repo *repository.Repository
}

// Login handles user authentication
// Steps: find user, verify password, generate token, update last login
func (s *authService) Login(ctx context.Context, data *dtos.LoginDTO) (*models.LoginResponse, error) {
	// Find user by email
	user, err := s.repo.User.FindByEmail(ctx, data.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Check if user is active
	if !user.IsActive {
		return nil, errors.New("account is deactivated")
	}

	// Verify password
	if !user.CheckPassword(data.Password) {
		return nil, errors.New("invalid credentials")
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Update last login timestamp
	if err := s.repo.User.UpdateLastLogin(ctx, user.ID.Hex()); err != nil {
		// Log error but don't fail the login
		// This is non-critical
		fmt.Printf("Failed to update last login: %v\n", err)
	}

	return &models.LoginResponse{
		User:  *user,
		Token: token,
	}, nil
}

func (s *authService) Register(ctx context.Context, data *dtos.RegisterDTO) (*models.User, error) {
	user := &models.User{
		FirstName: data.FirstName,
		LastName:  data.LastName,
		Email:     data.Email,
		Password:  data.Password,
		Role:      data.Role,
	}

	// Hash the password before storing
	if err := user.HashPassword(); err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Save user to database
	// This will fail if email already exists (unique index)
	if err := s.repo.User.Create(ctx, user); err != nil {
		return nil, err
	}

	// Return user without password (already excluded by JSON tag)
	return user, nil
}

// ValidateToken validates a JWT token and returns the user
func (s *authService) ValidateToken(tokenStr string) (*models.User, error) {
	// validate and parse token
	claims, err := utils.ValidateToken(tokenStr)
	if err != nil {
		return nil, err
	}

	userID, err := bson.ObjectIDFromHex(claims.ID)
	if err != nil {
		return nil, errors.New("invalid user ID in token")
	}

	// Create background context for DB operation
	ctx := context.Background()

	// Fetch user from database to ensure they still exist and are active
	user, err := s.repo.User.FindByID(ctx, userID.Hex())
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Check if user is still active
	if !user.IsActive {
		return nil, errors.New("account is deactivated")
	}

	return user, nil
}

// GetUserByID retrieves a user by their ID
func (s *authService) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	return s.repo.User.FindByID(ctx, id)
}

// CreateAdminUser creates an admin user if none exists
// This is for initial setup and should be called on application start
func (s *authService) CreateAdminUser(ctx context.Context, email, password string) error {
	// Check if admin already exists
	existing, _ := s.repo.User.FindByEmail(ctx, email)
	if existing != nil {
		return nil // Admin already exists
	}

	// Create admin user
	admin := &models.User{
		Email:    email,
		Password: password,
		Role:     "admin",
	}

	// Hash password
	if err := admin.HashPassword(); err != nil {
		return err
	}

	// Save to database
	return s.repo.User.Create(ctx, admin)
}
