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
	
	// Count items assigned to this participant
	itemCount := 0
	for _, assignment := range assignments {
		if assignment.ParticipantID == participant.ID {
			itemCount++
		}
	}
	
	body := fmt.Sprintf(`
Hello!

The assignment for group "%s" has been completed.

You have been assigned %d item(s).

To view the full results, please visit:
%s/participant/%s

Best regards,
OptiAssign Team
`, group.Name, itemCount, baseURL, participant.Token)

	if s.logToConsole {
		log.Printf("=== EMAIL NOTIFICATION ===")
		log.Printf("To: participant_%d", participant.ID)
		log.Printf("Subject: %s", subject)
		log.Printf("Body:\n%s", body)
		log.Printf("===========================")
	}

	return nil
}