package domain

import (
	"fmt"
	"log"
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
func (s *EmailService) SendInvitation(participant *Participant, group *Group, baseURL string) error {
	invitationURL := fmt.Sprintf("%s/participant/%s", baseURL, participant.Token)
	
	subject := fmt.Sprintf("Invitation to participate in: %s", group.Name)
	body := fmt.Sprintf(`
Hello!

You have been invited to participate in the group assignment: "%s"

Distribution Rule: %s

To submit your priorities, please visit:
%s

This link is unique to you and will allow you to rank the items in order of your preference.

Best regards,
OptiAssign Team
`, group.Name, group.Rule, invitationURL)

	if s.logToConsole {
		log.Printf("=== EMAIL INVITATION ===")
		log.Printf("To: participant_%d", participant.ID)
		log.Printf("Subject: %s", subject)
		log.Printf("Body:\n%s", body)
		log.Printf("========================")
	}

	// In a real implementation, this would send an actual email
	// For MVS, we just log to console
	
	return nil
}

// SendAssignmentNotification sends notification when assignment is complete
func (s *EmailService) SendAssignmentNotification(participant *Participant, group *Group, assignments []*Assignment, baseURL string) error {
	subject := fmt.Sprintf("Assignment Complete: %s", group.Name)
	
	// Get items assigned to this participant
	assignedItems := s.getAssignedItems(participant.ID, assignments)
	
	body := s.generateAssignmentEmailBody(participant, group, assignedItems, baseURL)

	if s.logToConsole {
		log.Printf("=== EMAIL NOTIFICATION ===")
		log.Printf("To: participant_%d", participant.ID)
		log.Printf("Subject: %s", subject)
		log.Printf("Body:\n%s", body)
		log.Printf("===========================")
	}

	return nil
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