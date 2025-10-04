package domain

import (
	"fmt"
	"log"
	"strings"
)

// EmailService handles email operations
type EmailService struct {
	logToConsole bool
}

// NewEmailService creates a new email service
func NewEmailService(logToConsole bool) *EmailService {
	return &EmailService{
		logToConsole: logToConsole,
	}
}

// SendInvitation sends an invitation email to a participant
func (s *EmailService) SendInvitation(participant *Participant, group *Group, items []*Item, baseURL string) error {
	invitationURL := fmt.Sprintf("%s/participant/%s", baseURL, participant.Token)
	
	subject := fmt.Sprintf("Invitation to participate in: %s", group.Name)
	
	// Create email template service
	templateService := NewEmailTemplateService(baseURL)
	
	// Prepare email data
	emailData := &InvitationEmailData{
		GroupName:        group.Name,
		ParticipantName:  fmt.Sprintf("Participant #%d", participant.ID),
		InvitationURL:    invitationURL,
		OwnerName:        "Group Owner", // In real app, get from user
		DistributionRule: group.Rule,
		ItemCount:        len(items),
		Items:            items,
	}
	
	// Generate HTML email
	htmlBody := templateService.GenerateInvitationHTML(emailData)
	
	// Generate plain text fallback
	textBody := fmt.Sprintf(`
Hello!

You have been invited to participate in the group assignment: "%s"

Distribution Rule: %s
Items to Rank: %d

Items in this assignment:
%s

To submit your priorities, please visit:
%s

This link is unique to you and will allow you to rank the items in order of your preference.

Best regards,
OptiAssign Team
`, group.Name, group.Rule, len(items), s.formatItemsList(items), invitationURL)

	if s.logToConsole {
		log.Printf("=== EMAIL INVITATION ===")
		log.Printf("To: participant_%d", participant.ID)
		log.Printf("Subject: %s", subject)
		log.Printf("HTML Body:\n%s", htmlBody)
		log.Printf("Text Body:\n%s", textBody)
		log.Printf("========================")
	}

	return nil
}

// formatItemsList formats items for plain text email
func (s *EmailService) formatItemsList(items []*Item) string {
	var list strings.Builder
	for i, item := range items {
		list.WriteString(fmt.Sprintf("%d. %s", i+1, item.Name))
		if item.Description != "" {
			list.WriteString(fmt.Sprintf(" - %s", item.Description))
		}
		list.WriteString("\n")
	}
	return list.String()
}

// SendAssignmentNotification sends notification when assignment is complete
func (s *EmailService) SendAssignmentNotification(participant *Participant, group *Group, assignments []*Assignment, items []*Item, baseURL string) error {
	subject := fmt.Sprintf("Your Assignment Results: %s", group.Name)
	
	// Get items assigned to this participant
	assignedItems := s.getAssignedItemsForParticipant(participant.ID, assignments, items)
	
	// Create email template service
	templateService := NewEmailTemplateService(baseURL)
	
	// Prepare email data
	emailData := &AssignmentEmailData{
		GroupName:        group.Name,
		ParticipantName:  fmt.Sprintf("Participant #%d", participant.ID),
		AssignedItems:    assignedItems,
		TotalItems:       len(items),
		ResultsURL:       fmt.Sprintf("%s/participant/%s/results", baseURL, participant.Token),
		OwnerName:        "Group Owner", // In real app, get from user
	}
	
	// Generate HTML email
	htmlBody := templateService.GenerateAssignmentHTML(emailData)
	
	// Generate plain text fallback
	textBody := s.generateAssignmentTextBody(participant, group, assignedItems, baseURL)

	if s.logToConsole {
		log.Printf("=== EMAIL NOTIFICATION ===")
		log.Printf("To: participant_%d", participant.ID)
		log.Printf("Subject: %s", subject)
		log.Printf("HTML Body:\n%s", htmlBody)
		log.Printf("Text Body:\n%s", textBody)
		log.Printf("===========================")
	}

	return nil
}

// getAssignedItemsForParticipant returns items assigned to a specific participant
func (s *EmailService) getAssignedItemsForParticipant(participantID int, assignments []*Assignment, items []*Item) []*Item {
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

// generateAssignmentTextBody creates plain text assignment email
func (s *EmailService) generateAssignmentTextBody(participant *Participant, group *Group, assignedItems []*Item, baseURL string) string {
	body := fmt.Sprintf(`
Hello!

The assignment for group "%s" has been completed.

You have been assigned %d item(s):

`, group.Name, len(assignedItems))

	// Add assigned items list
	for i, item := range assignedItems {
		body += fmt.Sprintf("%d. %s", i+1, item.Name)
		if item.Description != "" {
			body += fmt.Sprintf(" - %s", item.Description)
		}
		body += "\n"
	}

	body += fmt.Sprintf(`
To view the full results and details, please visit:
%s/participant/%s/results

Best regards,
OptiAssign Team
`, baseURL, participant.Token)

	return body
}

// getAssignedItems returns items assigned to a specific participant
func (s *EmailService) getAssignedItems(participantID int, assignments []*Assignment) []*Assignment {
	var assignedItems []*Assignment
	for _, assignment := range assignments {
		if assignment.ParticipantID == participantID {
			assignedItems = append(assignedItems, assignment)
		}
	}
	return assignedItems
}

// generateAssignmentEmailBody creates the email body for assignment notifications
func (s *EmailService) generateAssignmentEmailBody(participant *Participant, group *Group, assignedItems []*Assignment, baseURL string) string {
	itemCount := len(assignedItems)
	
	body := fmt.Sprintf(`
Hello!

The assignment for group "%s" has been completed.

You have been assigned %d item(s):

`, group.Name, itemCount)

	// Add assigned items list
	for i, assignment := range assignedItems {
		body += fmt.Sprintf("%d. Item ID: %d\n", i+1, assignment.ItemID)
	}

	body += fmt.Sprintf(`
To view the full results and details, please visit:
%s/participant/%s

Best regards,
OptiAssign Team
`, baseURL, participant.Token)

	return body
}