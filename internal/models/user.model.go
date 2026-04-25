// Package models
package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        bson.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName string        `bson:"firstName" json:"firstName"`
	LastName  string        `bson:"lastName" json:"lastName"`
	Email     string        `bson:"email" json:"email"`
	Password  string        `bson:"password" json:"-"`
	Role      string        `bson:"role" json:"role"`
	IsActive  bool          `bson:"isActive" json:"isActive"`
	CreatedAt time.Time     `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time     `bson:"updatedAt" json:"updatedAt"`
	LastLogin *time.Time    `bson:"lastLogin" json:"lastLogin"`
}

// HashPassword hashes the user's password using bcrypt
// bcrypt is a secure hashing algorithm designed for passwords
// It automatically handles salting and is computationally expensive
func (u *User) HashPassword() error {
	// GenerateFromPassword returns a hashed password
	// bcrypt.DefaultCost is 10 - good balance between security and performance
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return nil
}

// CheckPassword compares a plain text password with the hashed password
// Returns true if they match, false otherwise
func (u *User) CheckPassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) == nil
}

// PrepareUser sets up a new user before saving to database
// This ensures all required fields are properly initialized
func (u *User) PrepareUser() {
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()

	// Default role if not set
	if u.Role == "" {
		u.Role = "user"
	}

	// Always active by default
	u.IsActive = true
}

// UpdateLoginTime updates the last login timestamp
func (u *User) UpdateLoginTime() {
	now := time.Now()
	u.LastLogin = &now
	u.UpdatedAt = now
}
