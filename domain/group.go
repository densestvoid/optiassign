package domain

import (
	"time"
)

// Group represents the main assignment container
type Group struct {
	ID          int       `json:"id"`
	OwnerUserID int       `json:"owner_user_id"`
	Name        string    `json:"name"`
	Status      string    `json:"status"` // draft, prioritizing, completed
	Rule        string    `json:"rule"`   // equal, exhaust
	CreatedAt   time.Time `json:"created_at"`
}

// GroupRepository defines the interface for group data operations
type GroupRepository interface {
	Create(group *Group) error
	GetByID(id int) (*Group, error)
	GetByOwner(ownerUserID int) ([]*Group, error)
	Update(group *Group) error
	Delete(id int) error
}