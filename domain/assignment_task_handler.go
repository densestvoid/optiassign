package domain

import (
	"context"
	"fmt"
	"log"
)

// AssignmentTaskHandler handles assignment execution tasks
type AssignmentTaskHandler struct {
	assignmentAlgorithm *AssignmentAlgorithm
	emailService        *EmailService
	groupRepo           GroupRepository
	participantRepo     ParticipantRepository
	baseURL             string
}

// NewAssignmentTaskHandler creates a new assignment task handler
func NewAssignmentTaskHandler(
	assignmentAlgorithm *AssignmentAlgorithm,
	emailService *EmailService,
	groupRepo GroupRepository,
	participantRepo ParticipantRepository,
	baseURL string,
) *AssignmentTaskHandler {
	return &AssignmentTaskHandler{
		assignmentAlgorithm: assignmentAlgorithm,
		emailService:        emailService,
		groupRepo:           groupRepo,
		participantRepo:     participantRepo,
		baseURL:             baseURL,
	}
}

// Handle processes an assignment task
func (h *AssignmentTaskHandler) Handle(ctx context.Context, task Task) error {
	groupID, ok := task.Data["group_id"].(int)
	if !ok {
		return fmt.Errorf("invalid group_id in task data")
	}

	log.Printf("Executing assignment for group %d", groupID)

	// Execute assignment
	result, err := h.assignmentAlgorithm.ExecuteAssignment(ctx, groupID)
	if err != nil {
		return fmt.Errorf("assignment execution failed: %w", err)
	}

	// Get group details for email
	group, err := h.groupRepo.GetByID(groupID)
	if err != nil {
		log.Printf("Failed to get group %d for email: %v", groupID, err)
		return nil // Don't fail the task for email issues
	}

	// Get participants for email notifications
	participants, err := h.participantRepo.GetByGroupID(groupID)
	if err != nil {
		log.Printf("Failed to get participants for group %d: %v", groupID, err)
		return nil // Don't fail the task for email issues
	}

	// Send email notifications
	for _, participant := range participants {
		if err := h.emailService.SendAssignmentNotification(participant, group, result.Assignments, h.baseURL); err != nil {
			log.Printf("Failed to send email to participant %d: %v", participant.ID, err)
			// Continue with other participants
		}
	}

	log.Printf("Assignment completed for group %d with %d assignments", groupID, len(result.Assignments))
	return nil
}