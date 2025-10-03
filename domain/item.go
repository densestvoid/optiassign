package domain

import (
	"time"
)

// Item represents an assignment item
type Item struct {
	ID          int       `json:"id"`
	GroupID     int       `json:"group_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	IsAssigned  bool      `json:"is_assigned"`
	CreatedAt   time.Time `json:"created_at"`
}

// ItemRepository defines the interface for item data operations
type ItemRepository interface {
	Create(item *Item) error
	GetByGroupID(groupID int) ([]*Item, error)
	GetByID(id int) (*Item, error)
	Update(item *Item) error
	Delete(id int) error
	DeleteByGroupID(groupID int) error
}