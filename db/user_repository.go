package db

import (
	"database/sql"
	"fmt"
	"optiassign/domain"
)

type userRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository() domain.UserRepository {
	return &userRepository{db: DB}
}

func (r *userRepository) Create(user *domain.User) error {
	query := `
		INSERT INTO users (google_id, email, name) 
		VALUES ($1, $2, $3) 
		RETURNING id, created_at`
	
	err := r.db.QueryRow(query, user.GoogleID, user.Email, user.Name).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	
	return nil
}

func (r *userRepository) GetByGoogleID(googleID string) (*domain.User, error) {
	query := `SELECT id, google_id, email, name, created_at FROM users WHERE google_id = $1`
	
	user := &domain.User{}
	err := r.db.QueryRow(query, googleID).Scan(
		&user.ID, &user.GoogleID, &user.Email, &user.Name, &user.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by google_id: %w", err)
	}
	
	return user, nil
}

func (r *userRepository) GetByEmail(email string) (*domain.User, error) {
	query := `SELECT id, google_id, email, name, created_at FROM users WHERE email = $1`
	
	user := &domain.User{}
	err := r.db.QueryRow(query, email).Scan(
		&user.ID, &user.GoogleID, &user.Email, &user.Name, &user.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	
	return user, nil
}

func (r *userRepository) GetByID(id int) (*domain.User, error) {
	query := `SELECT id, google_id, email, name, created_at FROM users WHERE id = $1`
	
	user := &domain.User{}
	err := r.db.QueryRow(query, id).Scan(
		&user.ID, &user.GoogleID, &user.Email, &user.Name, &user.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}
	
	return user, nil
}

func (r *userRepository) Exists(email string) (bool, error) {
	query := `SELECT COUNT(*) FROM users WHERE email = $1`
	
	var count int
	err := r.db.QueryRow(query, email).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check user existence: %w", err)
	}
	
	return count > 0, nil
}