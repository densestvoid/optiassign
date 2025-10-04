package db

import (
	"database/sql"
	"fmt"
	"optiassign/domain"
)

type assignmentRepository struct {
	db *sql.DB
}

// NewAssignmentRepository creates a new assignment repository
func NewAssignmentRepository() domain.AssignmentRepository {
	return &assignmentRepository{db: DB}
}

func (r *assignmentRepository) Create(assignment *domain.Assignment) error {
	query := `
		INSERT INTO assignments (group_id, participant_id, item_id) 
		VALUES ($1, $2, $3) 
		RETURNING id, created_at`
	
	err := r.db.QueryRow(query, assignment.GroupID, assignment.ParticipantID, assignment.ItemID).Scan(&assignment.ID, &assignment.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create assignment: %w", err)
	}
	
	return nil
}

func (r *assignmentRepository) GetByGroupID(groupID int) ([]*domain.Assignment, error) {
	query := `SELECT id, group_id, participant_id, item_id, created_at FROM assignments WHERE group_id = $1 ORDER BY created_at ASC`
	
	rows, err := r.db.Query(query, groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to get assignments by group: %w", err)
	}
	defer rows.Close()
	
	var assignments []*domain.Assignment
	for rows.Next() {
		assignment := &domain.Assignment{}
		err := rows.Scan(
			&assignment.ID, &assignment.GroupID, &assignment.ParticipantID, &assignment.ItemID, &assignment.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan assignment: %w", err)
		}
		assignments = append(assignments, assignment)
	}
	
	return assignments, nil
}

func (r *assignmentRepository) GetByParticipantID(participantID int) ([]*domain.Assignment, error) {
	query := `SELECT id, group_id, participant_id, item_id, created_at FROM assignments WHERE participant_id = $1 ORDER BY created_at ASC`
	
	rows, err := r.db.Query(query, participantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get assignments by participant: %w", err)
	}
	defer rows.Close()
	
	var assignments []*domain.Assignment
	for rows.Next() {
		assignment := &domain.Assignment{}
		err := rows.Scan(
			&assignment.ID, &assignment.GroupID, &assignment.ParticipantID, &assignment.ItemID, &assignment.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan assignment: %w", err)
		}
		assignments = append(assignments, assignment)
	}
	
	return assignments, nil
}

func (r *assignmentRepository) DeleteByGroupID(groupID int) error {
	query := `DELETE FROM assignments WHERE group_id = $1`
	
	_, err := r.db.Exec(query, groupID)
	if err != nil {
		return fmt.Errorf("failed to delete assignments by group: %w", err)
	}
	
	return nil
}