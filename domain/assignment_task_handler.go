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
	itemRepo            ItemRepository
	baseURL             string
}

// NewAssignmentTaskHandler creates a new assignment task handler
func NewAssignmentTaskHandler(
	assignmentAlgorithm *AssignmentAlgorithm,
	emailService *EmailService,
	groupRepo GroupRepository,
	participantRepo ParticipantRepository,
	itemRepo ItemRepository,
	baseURL string,
) *AssignmentTaskHandler {
	return &AssignmentTaskHandler{
		assignmentAlgorithm: assignmentAlgorithm,
		emailService:        emailService,
		groupRepo:           groupRepo,
		participantRepo:     participantRepo,
		itemRepo:            itemRepo,
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

	// Get items for detailed email content
	items, err := h.itemRepo.GetByGroupID(groupID)
	if err != nil {
		log.Printf("Failed to get items for group %d: %v", groupID, err)
		return nil // Don't fail the task for email issues
	}

	// Send personalized email notifications with item details
	for _, participant := range participants {
		if err := h.sendPersonalizedAssignmentEmail(participant, group, result.Assignments, items); err != nil {
			log.Printf("Failed to send email to participant %d: %v", participant.ID, err)
			// Continue with other participants
		}
	}

	log.Printf("Assignment completed for group %d with %d assignments", groupID, len(result.Assignments))
	return nil
}

// sendPersonalizedAssignmentEmail sends a personalized email with assigned items
func (h *AssignmentTaskHandler) sendPersonalizedAssignmentEmail(participant *Participant, group *Group, assignments []*Assignment, items []*Item) error {
	// Get items assigned to this participant
	assignedItems := h.getAssignedItemsForParticipant(participant.ID, assignments, items)
	
	// Create detailed email content
	subject := fmt.Sprintf("Your Assignment Results: %s", group.Name)
	body := h.generateDetailedAssignmentEmail(participant, group, assignedItems)

	// Log the email (in production, this would send actual emails)
	log.Printf("=== PERSONALIZED EMAIL NOTIFICATION ===")
	log.Printf("To: participant_%d", participant.ID)
	log.Printf("Subject: %s", subject)
	log.Printf("Body:\n%s", body)
	log.Printf("=======================================")

	return nil
}

// getAssignedItemsForParticipant returns items assigned to a specific participant with details
func (h *AssignmentTaskHandler) getAssignedItemsForParticipant(participantID int, assignments []*Assignment, items []*Item) []*Item {
	// Create a map of item IDs for quick lookup
	itemMap := make(map[int]*Item)
	for _, item := range items {
		itemMap[item.ID] = item
	}

	var assignedItems []*Item
	for _, assignment := range assignments {
		if assignment.ParticipantID == participantID {
			if item, exists := itemMap[assignment.ItemID]; exists {
				assignedItems = append(assignedItems, item)
			}
		}
	}
	return assignedItems
}

// generateDetailedAssignmentEmail creates a detailed email with item information
func (h *AssignmentTaskHandler) generateDetailedAssignmentEmail(participant *Participant, group *Group, assignedItems []*Item) string {
	itemCount := len(assignedItems)
	
	body := fmt.Sprintf(`
Hello!

The assignment for group "%s" has been completed.

You have been assigned %d item(s):

`, group.Name, itemCount)

	// Add detailed item information
	for i, item := range assignedItems {
		body += fmt.Sprintf("%d. %s", i+1, item.Name)
		if item.Description != "" {
			body += fmt.Sprintf(" - %s", item.Description)
		}
		body += "\n"
	}

	body += fmt.Sprintf(`
To view the full results and details, please visit:
%s/participant/%s

Best regards,
OptiAssign Team
`, h.baseURL, participant.Token)

	return body
}