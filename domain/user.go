package domain

import (
	"time"
)

// User represents a registered user via Google SSO
type User struct {
	ID        int       `json:"id"`
	GoogleID  string    `json:"google_id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

// UserRepository defines the interface for user data operations
type UserRepository interface {
	Create(user *User) error
	GetByGoogleID(googleID string) (*User, error)
	GetByEmail(email string) (*User, error)
	GetByID(id int) (*User, error)
	Exists(email string) (bool, error)
}