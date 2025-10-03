package domain

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

// AssignmentAlgorithm implements the snaking draft assignment algorithm
type AssignmentAlgorithm struct {
	groupRepo        GroupRepository
	itemRepo         ItemRepository
	participantRepo  ParticipantRepository
	prioritizationRepo PrioritizationRepository
	assignmentRepo   AssignmentRepository
}

// NewAssignmentAlgorithm creates a new assignment algorithm
func NewAssignmentAlgorithm(
	groupRepo GroupRepository,
	itemRepo ItemRepository,
	participantRepo ParticipantRepository,
	prioritizationRepo PrioritizationRepository,
	assignmentRepo AssignmentRepository,
) *AssignmentAlgorithm {
	return &AssignmentAlgorithm{
		groupRepo:         groupRepo,
		itemRepo:          itemRepo,
		participantRepo:   participantRepo,
		prioritizationRepo: prioritizationRepo,
		assignmentRepo:    assignmentRepo,
	}
}

// AssignmentResult represents the result of an assignment
type AssignmentResult struct {
	GroupID      int                    `json:"group_id"`
	Assignments  []*Assignment          `json:"assignments"`
	Unassigned   []*Item                `json:"unassigned"`
	ParticipantOrder []int              `json:"participant_order"`
	Rounds       int                    `json:"rounds"`
}

// ExecuteAssignment runs the snaking draft assignment algorithm
func (a *AssignmentAlgorithm) ExecuteAssignment(ctx context.Context, groupID int) (*AssignmentResult, error) {
	// Get group
	group, err := a.groupRepo.GetByID(groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to get group: %w", err)
	}
	if group == nil {
		return nil, fmt.Errorf("group not found")
	}

	// Get items
	items, err := a.itemRepo.GetByGroupID(groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to get items: %w", err)
	}

	// Get participants
	participants, err := a.participantRepo.GetByGroupID(groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to get participants: %w", err)
	}

	if len(participants) == 0 {
		return nil, fmt.Errorf("no participants found")
	}

	if len(items) == 0 {
		return nil, fmt.Errorf("no items found")
	}

	// Check if all participants have submitted priorities
	for _, participant := range participants {
		hasSubmitted, err := a.prioritizationRepo.HasSubmitted(participant.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to check participant submission: %w", err)
		}
		if !hasSubmitted {
			return nil, fmt.Errorf("participant %d has not submitted priorities", participant.ID)
		}
	}

	// Clear existing assignments
	if err := a.assignmentRepo.DeleteByGroupID(groupID); err != nil {
		return nil, fmt.Errorf("failed to clear existing assignments: %w", err)
	}

	// Generate random participant order
	participantOrder := a.generateRandomOrder(participants)
	
	// Get all prioritizations
	allPrioritizations := make(map[int][]*Prioritization)
	for _, participant := range participants {
		prioritizations, err := a.prioritizationRepo.GetByParticipantID(participant.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get prioritizations: %w", err)
		}
		allPrioritizations[participant.ID] = prioritizations
	}

	// Execute snaking draft
	assignments, unassigned := a.executeSnakingDraft(items, participantOrder, allPrioritizations, group.Rule)

	// Save assignments
	for _, assignment := range assignments {
		if err := a.assignmentRepo.Create(assignment); err != nil {
			return nil, fmt.Errorf("failed to create assignment: %w", err)
		}
	}

	// Update group status
	group.Status = "completed"
	if err := a.groupRepo.Update(group); err != nil {
		return nil, fmt.Errorf("failed to update group status: %w", err)
	}

	// Calculate rounds
	rounds := len(items) / len(participants)
	if len(items)%len(participants) != 0 {
		rounds++
	}

	return &AssignmentResult{
		GroupID:          groupID,
		Assignments:      assignments,
		Unassigned:       unassigned,
		ParticipantOrder: participantOrder,
		Rounds:           rounds,
	}, nil
}

// generateRandomOrder creates a random order for participants
func (a *AssignmentAlgorithm) generateRandomOrder(participants []*Participant) []int {
	order := make([]int, len(participants))
	for i, participant := range participants {
		order[i] = participant.ID
	}
	
	// Shuffle the order
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(order), func(i, j int) {
		order[i], order[j] = order[j], order[i]
	})
	
	return order
}

// executeSnakingDraft runs the snaking draft algorithm
func (a *AssignmentAlgorithm) executeSnakingDraft(items []*Item, participantOrder []int, allPrioritizations map[int][]*Prioritization, rule string) ([]*Assignment, []*Item) {
	var assignments []*Assignment
	var unassigned []*Item
	
	// Create available items map
	availableItems := make(map[int]*Item)
	for _, item := range items {
		availableItems[item.ID] = item
	}
	
	// Calculate items per participant for equal distribution
	itemsPerParticipant := len(items) / len(participantOrder)
	remainingItems := len(items) % len(participantOrder)
	
	// Track items assigned per participant
	itemsAssigned := make(map[int]int)
	for _, participantID := range participantOrder {
		itemsAssigned[participantID] = 0
	}
	
	round := 0
	for len(availableItems) > 0 {
		// Determine pick order for this round (snaking)
		pickOrder := participantOrder
		if round%2 == 1 {
			// Reverse order for odd rounds
			pickOrder = make([]int, len(participantOrder))
			for i, participantID := range participantOrder {
				pickOrder[len(participantOrder)-1-i] = participantID
			}
		}
		
		// Each participant picks one item
		for _, participantID := range pickOrder {
			if len(availableItems) == 0 {
				break
			}
			
			// Check if participant has reached their limit (for equal distribution)
			if rule == "equal" && itemsAssigned[participantID] >= itemsPerParticipant {
				// Check if they can get one more (for remainder)
				if itemsAssigned[participantID] >= itemsPerParticipant+1 || remainingItems <= 0 {
					continue
				}
				remainingItems--
			}
			
			// Find best available item for this participant
			bestItem := a.findBestItem(participantID, availableItems, allPrioritizations)
			if bestItem == nil {
				continue
			}
			
			// Create assignment
			assignment := &Assignment{
				GroupID:       items[0].GroupID, // All items have same group ID
				ParticipantID: participantID,
				ItemID:        bestItem.ID,
			}
			assignments = append(assignments, assignment)
			
			// Remove item from available
			delete(availableItems, bestItem.ID)
			itemsAssigned[participantID]++
		}
		
		round++
	}
	
	// Mark remaining items as unassigned
	for _, item := range availableItems {
		unassigned = append(unassigned, item)
	}
	
	return assignments, unassigned
}

// findBestItem finds the best available item for a participant based on their priorities
func (a *AssignmentAlgorithm) findBestItem(participantID int, availableItems map[int]*Item, allPrioritizations map[int][]*Prioritization) *Item {
	prioritizations := allPrioritizations[participantID]
	if len(prioritizations) == 0 {
		return nil
	}
	
	// Sort prioritizations by rank (ascending)
	for i := 0; i < len(prioritizations)-1; i++ {
		for j := i + 1; j < len(prioritizations); j++ {
			if prioritizations[i].Rank > prioritizations[j].Rank {
				prioritizations[i], prioritizations[j] = prioritizations[j], prioritizations[i]
			}
		}
	}
	
	// Find first available item in priority order
	for _, prioritization := range prioritizations {
		if item, exists := availableItems[prioritization.ItemID]; exists {
			return item
		}
	}
	
	return nil
}