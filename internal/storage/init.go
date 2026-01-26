package storage

import (
	"time"

	"github.com/google/uuid"

	"driving-hours/internal/models"
)

// InitResult contains the result of initialization
type InitResult struct {
	AdminCreated  bool
	AdminEmail    string
	AdminPassword string
}

// PasswordHasher is a function type for hashing passwords
type PasswordHasher func(password string) (string, error)

// PasswordGenerator is a function type for generating random passwords
type PasswordGenerator func(length int) (string, error)

// Initialize checks for first run and creates admin if needed
func Initialize(storage Storage, hashPassword PasswordHasher, generatePassword PasswordGenerator) (*InitResult, error) {
	admin, err := storage.GetAdmin()
	if err != nil {
		return nil, err
	}

	// Admin already exists
	if admin != nil {
		return &InitResult{AdminCreated: false}, nil
	}

	// Generate random password
	password, err := generatePassword(16)
	if err != nil {
		return nil, err
	}

	// Hash the password
	hash, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	// Create admin user
	now := time.Now()
	admin = &models.User{
		ID:           uuid.New().String(),
		Email:        "admin@localhost",
		Name:         "Admin",
		PasswordHash: hash,
		Role:         models.RoleAdmin,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := storage.SaveAdmin(admin); err != nil {
		return nil, err
	}

	return &InitResult{
		AdminCreated:  true,
		AdminEmail:    admin.Email,
		AdminPassword: password,
	}, nil
}
