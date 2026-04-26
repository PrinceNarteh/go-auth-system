// Package utils
package utils

import (
	"errors"
	"time"

	"auth-system/internal/config"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// JWTClaims defines the structure of our JWT token payload
// These are the claims (data) stored in the token
type JWTClaims struct {
	UserID string `json:"userId"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateToken creates a new JWT token for a user
// Tokens contain user identity and expire after set time
func GenerateToken(userID bson.ObjectID, email, role string) (string, error) {
	// Calculate expiration time
	now := time.Now()
	expiresAt := now.Add(config.Env.JWT.Expiration)

	// Create claims with user data
	claims := &JWTClaims{
		UserID: userID.Hex(),
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			Issuer:    "auth-system",
			Subject:   userID.Hex(),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with our secret
	// This makes the token tamper-proof
	tokenString, err := token.SignedString(config.Env.JWT.Secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken verifies and parses a JWT token
// Returns the claims if valid, error otherwise
func ValidateToken(tokenStr string) (*JWTClaims, error) {
	// Parse and validate the token
	token, err := jwt.ParseWithClaims(tokenStr, &JWTClaims{}, func(t *jwt.Token) (any, error) {
		// Verify signing method
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(config.Env.JWT.Secret), nil // Verify signing method
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// RefreshToken generates a new token from an existing valid token
// This extends the session without requiring re-login
func RefreshToken(oldToken string) (string, error) {
	// Validate the old token first
	claims, err := ValidateToken(oldToken)
	if err != nil {
		return "", nil
	}

	userID, err := bson.ObjectIDFromHex(claims.ID)
	if err != nil {
		return "", nil
	}

	return GenerateToken(userID, claims.Email, claims.Role)
}
