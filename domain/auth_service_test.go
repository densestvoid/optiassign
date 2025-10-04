package domain

import (
	"context"
	"testing"
)

// MockUserRepository for testing
type MockUserRepository struct {
	users map[string]*User
}

func (m *MockUserRepository) Create(user *User) error {
	m.users[user.GoogleID] = user
	return nil
}

func (m *MockUserRepository) GetByGoogleID(googleID string) (*User, error) {
	if user, exists := m.users[googleID]; exists {
		return user, nil
	}
	return nil, nil
}

func (m *MockUserRepository) GetByEmail(email string) (*User, error) {
	for _, user := range m.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, nil
}

func (m *MockUserRepository) GetByID(id int) (*User, error) {
	for _, user := range m.users {
		if user.ID == id {
			return user, nil
		}
	}
	return nil, nil
}

func (m *MockUserRepository) Exists(email string) (bool, error) {
	for _, user := range m.users {
		if user.Email == email {
			return true, nil
		}
	}
	return false, nil
}

func TestAuthService_GetOrCreateUser(t *testing.T) {
	// Setup
	mockRepo := &MockUserRepository{
		users: make(map[string]*User),
	}
	
	// Add a pre-registered user
	preregisteredUser := &User{
		ID:       1,
		GoogleID: "",
		Email:    "test@example.com",
		Name:     "Test User",
	}
	mockRepo.users[""] = preregisteredUser
	
	authService := NewAuthService(mockRepo)
	
	// Test case: Get existing user by Google ID
	existingUser := &User{
		ID:       2,
		GoogleID: "google123",
		Email:    "existing@example.com",
		Name:     "Existing User",
	}
	mockRepo.users["google123"] = existingUser
	
	user, err := authService.GetOrCreateUser(context.Background(), "google123", "existing@example.com", "Existing User")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	if user.GoogleID != "google123" {
		t.Errorf("Expected GoogleID 'google123', got '%s'", user.GoogleID)
	}
	
	// Test case: Try to get non-pre-registered user
	_, err = authService.GetOrCreateUser(context.Background(), "google456", "notregistered@example.com", "Not Registered")
	if err == nil {
		t.Error("Expected error for non-pre-registered user, got nil")
	}
	
	// Test case: Get pre-registered user
	user, err = authService.GetOrCreateUser(context.Background(), "google789", "test@example.com", "Test User")
	if err != nil {
		t.Fatalf("Expected no error for pre-registered user, got %v", err)
	}
	
	if user.Email != "test@example.com" {
		t.Errorf("Expected email 'test@example.com', got '%s'", user.Email)
	}
}

func TestGenerateToken(t *testing.T) {
	token1, err := GenerateToken()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	if len(token1) == 0 {
		t.Error("Expected non-empty token")
	}
	
	// Generate another token to ensure they're different
	token2, err := GenerateToken()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	if token1 == token2 {
		t.Error("Expected different tokens, got same token")
	}
}