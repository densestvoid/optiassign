package db

import (
	"database/sql"
	"fmt"
	"optiassign/domain"
)

type participantRepository struct {
	db *sql.DB
}

// NewParticipantRepository creates a new participant repository
func NewParticipantRepository() domain.ParticipantRepository {
	return &participantRepository{db: DB}
}

func (r *participantRepository) Create(participant *domain.Participant) error {
	query := `
		INSERT INTO participants (group_id, user_id, token, has_submitted) 
		VALUES ($1, $2, $3, $4) 
		RETURNING id, created_at`
	
	err := r.db.QueryRow(query, participant.GroupID, participant.UserID, participant.Token, participant.HasSubmitted).Scan(&participant.ID, &participant.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create participant: %w", err)
	}
	
	return nil
}

func (r *participantRepository) GetByToken(token string) (*domain.Participant, error) {
	query := `SELECT id, group_id, user_id, token, has_submitted, created_at FROM participants WHERE token = $1`
	
	participant := &domain.Participant{}
	err := r.db.QueryRow(query, token).Scan(
		&participant.ID, &participant.GroupID, &participant.UserID, &participant.Token, &participant.HasSubmitted, &participant.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get participant by token: %w", err)
	}
	
	return participant, nil
}

func (r *participantRepository) GetByGroupID(groupID int) ([]*domain.Participant, error) {
	query := `SELECT id, group_id, user_id, token, has_submitted, created_at FROM participants WHERE group_id = $1 ORDER BY created_at ASC`
	
	rows, err := r.db.Query(query, groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to get participants by group: %w", err)
	}
	defer rows.Close()
	
	var participants []*domain.Participant
	for rows.Next() {
		participant := &domain.Participant{}
		err := rows.Scan(
			&participant.ID, &participant.GroupID, &participant.UserID, &participant.Token, &participant.HasSubmitted, &participant.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan participant: %w", err)
		}
		participants = append(participants, participant)
	}
	
	return participants, nil
}

func (r *participantRepository) GetByUserID(userID int) ([]*domain.Participant, error) {
	query := `SELECT id, group_id, user_id, token, has_submitted, created_at FROM participants WHERE user_id = $1 ORDER BY created_at ASC`
	
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get participants by user: %w", err)
	}
	defer rows.Close()
	
	var participants []*domain.Participant
	for rows.Next() {
		participant := &domain.Participant{}
		err := rows.Scan(
			&participant.ID, &participant.GroupID, &participant.UserID, &participant.Token, &participant.HasSubmitted, &participant.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan participant: %w", err)
		}
		participants = append(participants, participant)
	}
	
	return participants, nil
}

func (r *participantRepository) Update(participant *domain.Participant) error {
	query := `UPDATE participants SET has_submitted = $1 WHERE id = $2`
	
	_, err := r.db.Exec(query, participant.HasSubmitted, participant.ID)
	if err != nil {
		return fmt.Errorf("failed to update participant: %w", err)
	}
	
	return nil
}

func (r *participantRepository) Delete(id int) error {
	query := `DELETE FROM participants WHERE id = $1`
	
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete participant: %w", err)
	}
	
	return nil
}

func (r *participantRepository) DeleteByGroupID(groupID int) error {
	query := `DELETE FROM participants WHERE group_id = $1`
	
	_, err := r.db.Exec(query, groupID)
	if err != nil {
		return fmt.Errorf("failed to delete participants by group: %w", err)
	}
	
	return nil
}