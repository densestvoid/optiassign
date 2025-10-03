package domain

import (
	"time"
)

// Prioritization stores a user's ranked list of items
type Prioritization struct {
	ID            int       `json:"id"`
	ParticipantID int       `json:"participant_id"`
	ItemID        int       `json:"item_id"`
	Rank          int       `json:"rank"`
	CreatedAt     time.Time `json:"created_at"`
}

// PrioritizationRepository defines the interface for prioritization data operations
type PrioritizationRepository interface {
	Create(prioritization *Prioritization) error
	GetByParticipantID(participantID int) ([]*Prioritization, error)
	Update(prioritization *Prioritization) error
	DeleteByParticipantID(participantID int) error
	HasSubmitted(participantID int) (bool, error)
}