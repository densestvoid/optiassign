package db

import (
	"database/sql"
	"fmt"
	"optiassign/domain"
)

type itemRepository struct {
	db *sql.DB
}

// NewItemRepository creates a new item repository
func NewItemRepository() domain.ItemRepository {
	return &itemRepository{db: DB}
}

func (r *itemRepository) Create(item *domain.Item) error {
	query := `
		INSERT INTO items (group_id, name, description, is_assigned) 
		VALUES ($1, $2, $3, $4) 
		RETURNING id, created_at`
	
	err := r.db.QueryRow(query, item.GroupID, item.Name, item.Description, item.IsAssigned).Scan(&item.ID, &item.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create item: %w", err)
	}
	
	return nil
}

func (r *itemRepository) GetByGroupID(groupID int) ([]*domain.Item, error) {
	query := `SELECT id, group_id, name, description, is_assigned, created_at FROM items WHERE group_id = $1 ORDER BY created_at ASC`
	
	rows, err := r.db.Query(query, groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to get items by group: %w", err)
	}
	defer rows.Close()
	
	var items []*domain.Item
	for rows.Next() {
		item := &domain.Item{}
		err := rows.Scan(
			&item.ID, &item.GroupID, &item.Name, &item.Description, &item.IsAssigned, &item.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan item: %w", err)
		}
		items = append(items, item)
	}
	
	return items, nil
}

func (r *itemRepository) GetByID(id int) (*domain.Item, error) {
	query := `SELECT id, group_id, name, description, is_assigned, created_at FROM items WHERE id = $1`
	
	item := &domain.Item{}
	err := r.db.QueryRow(query, id).Scan(
		&item.ID, &item.GroupID, &item.Name, &item.Description, &item.IsAssigned, &item.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get item by id: %w", err)
	}
	
	return item, nil
}

func (r *itemRepository) Update(item *domain.Item) error {
	query := `UPDATE items SET name = $1, description = $2, is_assigned = $3 WHERE id = $4`
	
	_, err := r.db.Exec(query, item.Name, item.Description, item.IsAssigned, item.ID)
	if err != nil {
		return fmt.Errorf("failed to update item: %w", err)
	}
	
	return nil
}

func (r *itemRepository) Delete(id int) error {
	query := `DELETE FROM items WHERE id = $1`
	
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete item: %w", err)
	}
	
	return nil
}

func (r *itemRepository) DeleteByGroupID(groupID int) error {
	query := `DELETE FROM items WHERE group_id = $1`
	
	_, err := r.db.Exec(query, groupID)
	if err != nil {
		return fmt.Errorf("failed to delete items by group: %w", err)
	}
	
	return nil
}