package db

import (
	"database/sql"
	"fmt"
	"optiassign/domain"
)

type prioritizationRepository struct {
	db *sql.DB
}

// NewPrioritizationRepository creates a new prioritization repository
func NewPrioritizationRepository() domain.PrioritizationRepository {
	return &prioritizationRepository{db: DB}
}

func (r *prioritizationRepository) Create(prioritization *domain.Prioritization) error {
	query := `
		INSERT INTO prioritizations (participant_id, item_id, rank) 
		VALUES ($1, $2, $3) 
		RETURNING id, created_at`
	
	err := r.db.QueryRow(query, prioritization.ParticipantID, prioritization.ItemID, prioritization.Rank).Scan(&prioritization.ID, &prioritization.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create prioritization: %w", err)
	}
	
	return nil
}

func (r *prioritizationRepository) GetByParticipantID(participantID int) ([]*domain.Prioritization, error) {
	query := `SELECT id, participant_id, item_id, rank, created_at FROM prioritizations WHERE participant_id = $1 ORDER BY rank ASC`
	
	rows, err := r.db.Query(query, participantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get prioritizations by participant: %w", err)
	}
	defer rows.Close()
	
	var prioritizations []*domain.Prioritization
	for rows.Next() {
		prioritization := &domain.Prioritization{}
		err := rows.Scan(
			&prioritization.ID, &prioritization.ParticipantID, &prioritization.ItemID, &prioritization.Rank, &prioritization.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan prioritization: %w", err)
		}
		prioritizations = append(prioritizations, prioritization)
	}
	
	return prioritizations, nil
}

func (r *prioritizationRepository) Update(prioritization *domain.Prioritization) error {
	query := `UPDATE prioritizations SET rank = $1 WHERE id = $2`
	
	_, err := r.db.Exec(query, prioritization.Rank, prioritization.ID)
	if err != nil {
		return fmt.Errorf("failed to update prioritization: %w", err)
	}
	
	return nil
}

func (r *prioritizationRepository) DeleteByParticipantID(participantID int) error {
	query := `DELETE FROM prioritizations WHERE participant_id = $1`
	
	_, err := r.db.Exec(query, participantID)
	if err != nil {
		return fmt.Errorf("failed to delete prioritizations by participant: %w", err)
	}
	
	return nil
}

func (r *prioritizationRepository) HasSubmitted(participantID int) (bool, error) {
	query := `SELECT COUNT(*) FROM prioritizations WHERE participant_id = $1`
	
	var count int
	err := r.db.QueryRow(query, participantID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check prioritization submission: %w", err)
	}
	
	return count > 0, nil
}