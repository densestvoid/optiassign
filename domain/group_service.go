package domain

import (
	"context"
	"fmt"
)

// GroupService handles group business logic
type GroupService struct {
	groupRepo       GroupRepository
	itemRepo        ItemRepository
	participantRepo ParticipantRepository
	userRepo        UserRepository
}

// NewGroupService creates a new group service
func NewGroupService(
	groupRepo GroupRepository,
	itemRepo ItemRepository,
	participantRepo ParticipantRepository,
	userRepo UserRepository,
) *GroupService {
	return &GroupService{
		groupRepo:       groupRepo,
		itemRepo:        itemRepo,
		participantRepo: participantRepo,
		userRepo:        userRepo,
	}
}

// CreateGroupWithItems creates a group with items
func (s *GroupService) CreateGroupWithItems(ctx context.Context, ownerUserID int, groupName, rule string, itemNames []string) (*Group, error) {
	// Validate inputs
	if groupName == "" {
		return nil, fmt.Errorf("group name is required")
	}
	if rule != "equal" && rule != "exhaust" {
		return nil, fmt.Errorf("rule must be 'equal' or 'exhaust'")
	}
	if len(itemNames) == 0 {
		return nil, fmt.Errorf("at least one item is required")
	}

	// Create group
	group := &Group{
		OwnerUserID: ownerUserID,
		Name:        groupName,
		Status:      "draft",
		Rule:        rule,
	}

	if err := s.groupRepo.Create(group); err != nil {
		return nil, fmt.Errorf("failed to create group: %w", err)
	}

	// Create items
	for _, itemName := range itemNames {
		if itemName == "" {
			continue // Skip empty item names
		}
		
		item := &Item{
			GroupID:     group.ID,
			Name:        itemName,
			Description: "",
			IsAssigned:  false,
		}
		
		if err := s.itemRepo.Create(item); err != nil {
			return nil, fmt.Errorf("failed to create item: %w", err)
		}
	}

	return group, nil
}

// AddParticipant adds a participant to a group
func (s *GroupService) AddParticipant(ctx context.Context, groupID, userID int) (*Participant, error) {
	// Check if group exists
	group, err := s.groupRepo.GetByID(groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to get group: %w", err)
	}
	if group == nil {
		return nil, fmt.Errorf("group not found")
	}

	// Check if user exists
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	// Check if user is already a participant
	existingParticipants, err := s.participantRepo.GetByGroupID(groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing participants: %w", err)
	}
	
	for _, p := range existingParticipants {
		if p.UserID == userID {
			return nil, fmt.Errorf("user is already a participant")
		}
	}

	// Generate token
	token, err := GenerateToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Create participant
	participant := &Participant{
		GroupID:      groupID,
		UserID:       userID,
		Token:        token,
		HasSubmitted: false,
	}

	if err := s.participantRepo.Create(participant); err != nil {
		return nil, fmt.Errorf("failed to create participant: %w", err)
	}

	return participant, nil
}

// GetGroupWithDetails gets a group with its items and participants
func (s *GroupService) GetGroupWithDetails(ctx context.Context, groupID int) (*Group, []*Item, []*Participant, error) {
	// Get group
	group, err := s.groupRepo.GetByID(groupID)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get group: %w", err)
	}
	if group == nil {
		return nil, nil, nil, fmt.Errorf("group not found")
	}

	// Get items
	items, err := s.itemRepo.GetByGroupID(groupID)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get items: %w", err)
	}

	// Get participants
	participants, err := s.participantRepo.GetByGroupID(groupID)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get participants: %w", err)
	}

	return group, items, participants, nil
}

// UpdateGroupStatus updates the group status
func (s *GroupService) UpdateGroupStatus(ctx context.Context, groupID int, status string) error {
	group, err := s.groupRepo.GetByID(groupID)
	if err != nil {
		return fmt.Errorf("failed to get group: %w", err)
	}
	if group == nil {
		return fmt.Errorf("group not found")
	}

	group.Status = status
	return s.groupRepo.Update(group)
}

// GetUserGroups gets all groups for a user
func (s *GroupService) GetUserGroups(ctx context.Context, userID int) ([]*Group, error) {
	return s.groupRepo.GetByOwner(userID)
}

// DeleteGroup deletes a group and all its associated data
func (s *GroupService) DeleteGroup(ctx context.Context, groupID int) error {
	// Delete items
	if err := s.itemRepo.DeleteByGroupID(groupID); err != nil {
		return fmt.Errorf("failed to delete items: %w", err)
	}

	// Delete participants
	if err := s.participantRepo.DeleteByGroupID(groupID); err != nil {
		return fmt.Errorf("failed to delete participants: %w", err)
	}

	// Delete group
	return s.groupRepo.Delete(groupID)
}