package db

import (
	"database/sql"
	"fmt"
	"optiassign/domain"
)

type groupRepository struct {
	db *sql.DB
}

// NewGroupRepository creates a new group repository
func NewGroupRepository() domain.GroupRepository {
	return &groupRepository{db: DB}
}

func (r *groupRepository) Create(group *domain.Group) error {
	query := `
		INSERT INTO groups (owner_user_id, name, status, rule) 
		VALUES ($1, $2, $3, $4) 
		RETURNING id, created_at`
	
	err := r.db.QueryRow(query, group.OwnerUserID, group.Name, group.Status, group.Rule).Scan(&group.ID, &group.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create group: %w", err)
	}
	
	return nil
}

func (r *groupRepository) GetByID(id int) (*domain.Group, error) {
	query := `SELECT id, owner_user_id, name, status, rule, created_at FROM groups WHERE id = $1`
	
	group := &domain.Group{}
	err := r.db.QueryRow(query, id).Scan(
		&group.ID, &group.OwnerUserID, &group.Name, &group.Status, &group.Rule, &group.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get group by id: %w", err)
	}
	
	return group, nil
}

func (r *groupRepository) GetByOwner(ownerUserID int) ([]*domain.Group, error) {
	query := `SELECT id, owner_user_id, name, status, rule, created_at FROM groups WHERE owner_user_id = $1 ORDER BY created_at DESC`
	
	rows, err := r.db.Query(query, ownerUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get groups by owner: %w", err)
	}
	defer rows.Close()
	
	var groups []*domain.Group
	for rows.Next() {
		group := &domain.Group{}
		err := rows.Scan(
			&group.ID, &group.OwnerUserID, &group.Name, &group.Status, &group.Rule, &group.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan group: %w", err)
		}
		groups = append(groups, group)
	}
	
	return groups, nil
}

func (r *groupRepository) Update(group *domain.Group) error {
	query := `UPDATE groups SET name = $1, status = $2, rule = $3 WHERE id = $4`
	
	_, err := r.db.Exec(query, group.Name, group.Status, group.Rule, group.ID)
	if err != nil {
		return fmt.Errorf("failed to update group: %w", err)
	}
	
	return nil
}

func (r *groupRepository) Delete(id int) error {
	query := `DELETE FROM groups WHERE id = $1`
	
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete group: %w", err)
	}
	
	return nil
}