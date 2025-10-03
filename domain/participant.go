package domain

import (
	"time"
)

// Participant links a registered user to a group
type Participant struct {
	ID            int       `json:"id"`
	GroupID       int       `json:"group_id"`
	UserID        int       `json:"user_id"`
	Token         string    `json:"token"`
	HasSubmitted  bool      `json:"has_submitted"`
	CreatedAt     time.Time `json:"created_at"`
}

// ParticipantRepository defines the interface for participant data operations
type ParticipantRepository interface {
	Create(participant *Participant) error
	GetByToken(token string) (*Participant, error)
	GetByGroupID(groupID int) ([]*Participant, error)
	GetByUserID(userID int) ([]*Participant, error)
	Update(participant *Participant) error
	Delete(id int) error
	DeleteByGroupID(groupID int) error
}