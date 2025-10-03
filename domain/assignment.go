package domain

import (
	"time"
)

// Assignment represents the final result of the draft
type Assignment struct {
	ID            int       `json:"id"`
	GroupID       int       `json:"group_id"`
	ParticipantID int       `json:"participant_id"`
	ItemID        int       `json:"item_id"`
	CreatedAt     time.Time `json:"created_at"`
}

// AssignmentRepository defines the interface for assignment data operations
type AssignmentRepository interface {
	Create(assignment *Assignment) error
	GetByGroupID(groupID int) ([]*Assignment, error)
	GetByParticipantID(participantID int) ([]*Assignment, error)
	DeleteByGroupID(groupID int) error
}