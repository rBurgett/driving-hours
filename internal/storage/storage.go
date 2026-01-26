package storage

import (
	"driving-hours/internal/models"
)

// Storage defines the interface for data persistence
type Storage interface {
	// User operations
	GetUser(id string) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	GetAllUsers() ([]*models.User, error)
	GetDrivers() ([]*models.User, error)
	SaveUser(user *models.User) error
	DeleteUser(id string) error

	// Session operations
	GetSession(token string) (*models.Session, error)
	SaveSession(session *models.Session) error
	DeleteSession(token string) error
	CleanExpiredSessions() error

	// Admin operations
	GetAdmin() (*models.User, error)
	SaveAdmin(admin *models.User) error
}
