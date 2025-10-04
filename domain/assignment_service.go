package domain

import (
	"context"
	"fmt"
	"log"
	"time"
)

// AssignmentService manages assignment execution and completion checking
type AssignmentService struct {
	groupRepo           GroupRepository
	participantRepo     ParticipantRepository
	prioritizationRepo  PrioritizationRepository
	taskManager         *TaskManager
}

// NewAssignmentService creates a new assignment service
func NewAssignmentService(
	groupRepo GroupRepository,
	participantRepo ParticipantRepository,
	prioritizationRepo PrioritizationRepository,
	taskManager *TaskManager,
) *AssignmentService {
	return &AssignmentService{
		groupRepo:          groupRepo,
		participantRepo:    participantRepo,
		prioritizationRepo: prioritizationRepo,
		taskManager:        taskManager,
	}
}

// CheckAndTriggerAssignment checks if all participants have submitted and triggers assignment if ready
func (s *AssignmentService) CheckAndTriggerAssignment(ctx context.Context, groupID int) error {
	// Get group
	group, err := s.groupRepo.GetByID(groupID)
	if err != nil {
		return fmt.Errorf("failed to get group: %w", err)
	}
	if group == nil {
		return fmt.Errorf("group not found")
	}

	// Check if group is already completed
	if group.Status == "completed" {
		log.Printf("Group %d is already completed", groupID)
		return nil
	}

	// Get all participants
	participants, err := s.participantRepo.GetByGroupID(groupID)
	if err != nil {
		return fmt.Errorf("failed to get participants: %w", err)
	}

	if len(participants) == 0 {
		return fmt.Errorf("no participants found for group %d", groupID)
	}

	// Check if all participants have submitted
	allSubmitted := true
	for _, participant := range participants {
		hasSubmitted, err := s.prioritizationRepo.HasSubmitted(participant.ID)
		if err != nil {
			return fmt.Errorf("failed to check participant %d submission: %w", participant.ID, err)
		}
		if !hasSubmitted {
			allSubmitted = false
			log.Printf("Participant %d has not submitted priorities yet", participant.ID)
			break
		}
	}

	if !allSubmitted {
		log.Printf("Not all participants have submitted for group %d", groupID)
		return nil
	}

	// All participants have submitted - trigger assignment
	log.Printf("All participants have submitted for group %d, triggering assignment", groupID)

	// Create assignment task
	task := Task{
		ID:        fmt.Sprintf("assignment_%d_%d", groupID, time.Now().Unix()),
		Type:      "assignment",
		Data:      map[string]interface{}{"group_id": groupID},
		CreatedAt: time.Now(),
	}

	// Enqueue the task
	if err := s.taskManager.Enqueue(task); err != nil {
		return fmt.Errorf("failed to enqueue assignment task: %w", err)
	}

	log.Printf("Assignment task enqueued for group %d", groupID)
	return nil
}

// GetAssignmentStatus returns the status of assignment for a group
func (s *AssignmentService) GetAssignmentStatus(ctx context.Context, groupID int) (*AssignmentStatus, error) {
	// Get group
	group, err := s.groupRepo.GetByID(groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to get group: %w", err)
	}
	if group == nil {
		return nil, fmt.Errorf("group not found")
	}

	// Get participants
	participants, err := s.participantRepo.GetByGroupID(groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to get participants: %w", err)
	}

	// Check submission status
	submittedCount := 0
	totalCount := len(participants)
	
	for _, participant := range participants {
		hasSubmitted, err := s.prioritizationRepo.HasSubmitted(participant.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to check participant %d submission: %w", participant.ID, err)
		}
		if hasSubmitted {
			submittedCount++
		}
	}

	status := &AssignmentStatus{
		GroupID:        groupID,
		GroupStatus:    group.Status,
		TotalParticipants: totalCount,
		SubmittedCount:   submittedCount,
		AllSubmitted:     submittedCount == totalCount,
		IsCompleted:      group.Status == "completed",
	}

	return status, nil
}

// AssignmentStatus represents the status of an assignment
type AssignmentStatus struct {
	GroupID           int    `json:"group_id"`
	GroupStatus       string `json:"group_status"`
	TotalParticipants int    `json:"total_participants"`
	SubmittedCount    int    `json:"submitted_count"`
	AllSubmitted      bool   `json:"all_submitted"`
	IsCompleted       bool   `json:"is_completed"`
}