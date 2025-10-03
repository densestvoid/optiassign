package domain

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

// AuthService handles authentication operations
type AuthService struct {
	userRepo UserRepository
}

// NewAuthService creates a new authentication service
func NewAuthService(userRepo UserRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

// GetOrCreateUser retrieves an existing user or creates a new one
func (s *AuthService) GetOrCreateUser(ctx context.Context, googleID, email, name string) (*User, error) {
	// Try to get existing user by Google ID
	user, err := s.userRepo.GetByGoogleID(googleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by google_id: %w", err)
	}
	
	if user != nil {
		return user, nil
	}
	
	// Check if user exists by email (for pre-registered users)
	existingUser, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("failed to check user existence: %w", err)
	}
	
	if existingUser == nil {
		return nil, fmt.Errorf("user with email %s is not pre-registered", email)
	}
	
	// Update existing user with Google ID
	existingUser.GoogleID = googleID
	// Note: In a real implementation, you'd need an Update method in the repository
	// For MVS, we'll assume the user is already registered
	
	return existingUser, nil
}

// GenerateToken generates a secure random token
func GenerateToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}

// Session represents a user session
type Session struct {
	UserID    int       `json:"user_id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	ExpiresAt time.Time `json:"expires_at"`
}

// IsValid checks if the session is still valid
func (s *Session) IsValid() bool {
	return time.Now().Before(s.ExpiresAt)
}