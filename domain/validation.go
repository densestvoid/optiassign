package domain

import (
	"fmt"
	"regexp"
	"strings"
)

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// ValidationErrors represents multiple validation errors
type ValidationErrors struct {
	Errors []ValidationError `json:"errors"`
}

func (e ValidationErrors) Error() string {
	var messages []string
	for _, err := range e.Errors {
		messages = append(messages, err.Error())
	}
	return strings.Join(messages, "; ")
}

// HasErrors returns true if there are validation errors
func (e ValidationErrors) HasErrors() bool {
	return len(e.Errors) > 0
}

// AddError adds a validation error
func (e *ValidationErrors) AddError(field, message string) {
	e.Errors = append(e.Errors, ValidationError{Field: field, Message: message})
}

// ValidateGroup validates a group
func ValidateGroup(group *Group) *ValidationErrors {
	errors := &ValidationErrors{}
	
	if strings.TrimSpace(group.Name) == "" {
		errors.AddError("name", "Group name is required")
	} else if len(group.Name) > 100 {
		errors.AddError("name", "Group name must be 100 characters or less")
	}
	
	if group.Rule != "equal" && group.Rule != "uneven" {
		errors.AddError("rule", "Distribution rule must be 'equal' or 'uneven'")
	}
	
	return errors
}

// ValidateItem validates an item
func ValidateItem(item *Item) *ValidationErrors {
	errors := &ValidationErrors{}
	
	if strings.TrimSpace(item.Name) == "" {
		errors.AddError("name", "Item name is required")
	} else if len(item.Name) > 200 {
		errors.AddError("name", "Item name must be 200 characters or less")
	}
	
	if len(item.Description) > 500 {
		errors.AddError("description", "Item description must be 500 characters or less")
	}
	
	return errors
}

// ValidatePriorities validates priority submissions
func ValidatePriorities(priorities map[int]int, totalItems int) *ValidationErrors {
	errors := &ValidationErrors{}
	
	if len(priorities) != totalItems {
		errors.AddError("priorities", "All items must be ranked")
		return errors
	}
	
	// Check for duplicate ranks
	usedRanks := make(map[int]bool)
	for itemID, rank := range priorities {
		if rank < 1 || rank > totalItems {
			errors.AddError(fmt.Sprintf("item_%d", itemID), "Rank must be between 1 and "+fmt.Sprintf("%d", totalItems))
			continue
		}
		
		if usedRanks[rank] {
			errors.AddError(fmt.Sprintf("item_%d", itemID), "Each rank can only be used once")
		}
		usedRanks[rank] = true
	}
	
	return errors
}

// ValidateEmail validates email format
func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// ValidateToken validates participant token format
func ValidateToken(token string) bool {
	// Token should be alphanumeric and 32 characters long
	tokenRegex := regexp.MustCompile(`^[a-zA-Z0-9]{32}$`)
	return tokenRegex.MatchString(token)
}