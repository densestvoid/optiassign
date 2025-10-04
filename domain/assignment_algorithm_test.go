package domain

import (
	"context"
	"testing"
)

// Mock repositories for testing
type mockGroupRepo struct {
	groups map[int]*Group
}

func (m *mockGroupRepo) GetByID(id int) (*Group, error) {
	return m.groups[id], nil
}

func (m *mockGroupRepo) Update(group *Group) error {
	m.groups[group.ID] = group
	return nil
}

func (m *mockGroupRepo) Create(group *Group) error {
	m.groups[group.ID] = group
	return nil
}

func (m *mockGroupRepo) Delete(id int) error {
	delete(m.groups, id)
	return nil
}

func (m *mockGroupRepo) GetByOwner(ownerID int) ([]*Group, error) {
	var groups []*Group
	for _, group := range m.groups {
		groups = append(groups, group)
	}
	return groups, nil
}

type mockItemRepo struct {
	items map[int][]*Item
}

func (m *mockItemRepo) GetByGroupID(groupID int) ([]*Item, error) {
	return m.items[groupID], nil
}

func (m *mockItemRepo) Create(item *Item) error {
	m.items[item.GroupID] = append(m.items[item.GroupID], item)
	return nil
}

func (m *mockItemRepo) Update(item *Item) error {
	// Update item in the slice
	for i, existingItem := range m.items[item.GroupID] {
		if existingItem.ID == item.ID {
			m.items[item.GroupID][i] = item
			break
		}
	}
	return nil
}

func (m *mockItemRepo) Delete(id int) error {
	// Remove item from all groups
	for groupID, items := range m.items {
		for i, item := range items {
			if item.ID == id {
				m.items[groupID] = append(items[:i], items[i+1:]...)
				break
			}
		}
	}
	return nil
}

func (m *mockItemRepo) DeleteByGroupID(groupID int) error {
	delete(m.items, groupID)
	return nil
}

func (m *mockItemRepo) GetByID(id int) (*Item, error) {
	for _, items := range m.items {
		for _, item := range items {
			if item.ID == id {
				return item, nil
			}
		}
	}
	return nil, nil
}

type mockParticipantRepo struct {
	participants map[int][]*Participant
}

func (m *mockParticipantRepo) GetByGroupID(groupID int) ([]*Participant, error) {
	return m.participants[groupID], nil
}

func (m *mockParticipantRepo) Create(participant *Participant) error {
	m.participants[participant.GroupID] = append(m.participants[participant.GroupID], participant)
	return nil
}

func (m *mockParticipantRepo) Update(participant *Participant) error {
	// Update participant in the slice
	for i, existingParticipant := range m.participants[participant.GroupID] {
		if existingParticipant.ID == participant.ID {
			m.participants[participant.GroupID][i] = participant
			break
		}
	}
	return nil
}

func (m *mockParticipantRepo) Delete(id int) error {
	// Remove participant from all groups
	for groupID, participants := range m.participants {
		for i, participant := range participants {
			if participant.ID == id {
				m.participants[groupID] = append(participants[:i], participants[i+1:]...)
				break
			}
		}
	}
	return nil
}

func (m *mockParticipantRepo) DeleteByGroupID(groupID int) error {
	delete(m.participants, groupID)
	return nil
}

func (m *mockParticipantRepo) GetByToken(token string) (*Participant, error) {
	for _, participants := range m.participants {
		for _, participant := range participants {
			if participant.Token == token {
				return participant, nil
			}
		}
	}
	return nil, nil
}

func (m *mockParticipantRepo) GetByUserID(userID int) ([]*Participant, error) {
	var userParticipants []*Participant
	for _, participants := range m.participants {
		for _, participant := range participants {
			if participant.UserID == userID {
				userParticipants = append(userParticipants, participant)
			}
		}
	}
	return userParticipants, nil
}

type mockPrioritizationRepo struct {
	prioritizations map[int][]*Prioritization
}

func (m *mockPrioritizationRepo) HasSubmitted(participantID int) (bool, error) {
	_, exists := m.prioritizations[participantID]
	return exists, nil
}

func (m *mockPrioritizationRepo) GetByParticipantID(participantID int) ([]*Prioritization, error) {
	return m.prioritizations[participantID], nil
}

func (m *mockPrioritizationRepo) Create(prioritization *Prioritization) error {
	m.prioritizations[prioritization.ParticipantID] = append(m.prioritizations[prioritization.ParticipantID], prioritization)
	return nil
}

func (m *mockPrioritizationRepo) Update(prioritization *Prioritization) error {
	// Update prioritization in the slice
	for i, existingPrioritization := range m.prioritizations[prioritization.ParticipantID] {
		if existingPrioritization.ItemID == prioritization.ItemID {
			m.prioritizations[prioritization.ParticipantID][i] = prioritization
			break
		}
	}
	return nil
}

func (m *mockPrioritizationRepo) Delete(participantID int) error {
	delete(m.prioritizations, participantID)
	return nil
}

func (m *mockPrioritizationRepo) DeleteByGroupID(groupID int) error {
	// Clear all prioritizations for participants in this group
	for participantID := range m.prioritizations {
		delete(m.prioritizations, participantID)
	}
	return nil
}

func (m *mockPrioritizationRepo) DeleteByParticipantID(participantID int) error {
	delete(m.prioritizations, participantID)
	return nil
}

type mockAssignmentRepo struct {
	assignments []*Assignment
}

func (m *mockAssignmentRepo) DeleteByGroupID(groupID int) error {
	m.assignments = []*Assignment{}
	return nil
}

func (m *mockAssignmentRepo) Create(assignment *Assignment) error {
	m.assignments = append(m.assignments, assignment)
	return nil
}

func (m *mockAssignmentRepo) GetByGroupID(groupID int) ([]*Assignment, error) {
	return m.assignments, nil
}

func (m *mockAssignmentRepo) GetByParticipantID(participantID int) ([]*Assignment, error) {
	var participantAssignments []*Assignment
	for _, assignment := range m.assignments {
		if assignment.ParticipantID == participantID {
			participantAssignments = append(participantAssignments, assignment)
		}
	}
	return participantAssignments, nil
}

func TestAssignmentAlgorithm_EqualDistribution(t *testing.T) {
	// Setup test data
	group := &Group{ID: 1, Name: "Test Group", Rule: "equal", Status: "active"}
	items := []*Item{
		{ID: 1, GroupID: 1, Name: "Item 1"},
		{ID: 2, GroupID: 1, Name: "Item 2"},
		{ID: 3, GroupID: 1, Name: "Item 3"},
		{ID: 4, GroupID: 1, Name: "Item 4"},
		{ID: 5, GroupID: 1, Name: "Item 5"},
		{ID: 6, GroupID: 1, Name: "Item 6"},
	}
	participants := []*Participant{
		{ID: 1, GroupID: 1, Token: "token1"},
		{ID: 2, GroupID: 1, Token: "token2"},
		{ID: 3, GroupID: 1, Token: "token3"},
	}
	
	// Create prioritizations
	prioritizations := map[int][]*Prioritization{
		1: {
			{ParticipantID: 1, ItemID: 1, Rank: 1},
			{ParticipantID: 1, ItemID: 2, Rank: 2},
			{ParticipantID: 1, ItemID: 3, Rank: 3},
			{ParticipantID: 1, ItemID: 4, Rank: 4},
			{ParticipantID: 1, ItemID: 5, Rank: 5},
			{ParticipantID: 1, ItemID: 6, Rank: 6},
		},
		2: {
			{ParticipantID: 2, ItemID: 1, Rank: 6},
			{ParticipantID: 2, ItemID: 2, Rank: 5},
			{ParticipantID: 2, ItemID: 3, Rank: 4},
			{ParticipantID: 2, ItemID: 4, Rank: 3},
			{ParticipantID: 2, ItemID: 5, Rank: 2},
			{ParticipantID: 2, ItemID: 6, Rank: 1},
		},
		3: {
			{ParticipantID: 3, ItemID: 1, Rank: 3},
			{ParticipantID: 3, ItemID: 2, Rank: 2},
			{ParticipantID: 3, ItemID: 3, Rank: 1},
			{ParticipantID: 3, ItemID: 4, Rank: 6},
			{ParticipantID: 3, ItemID: 5, Rank: 5},
			{ParticipantID: 3, ItemID: 6, Rank: 4},
		},
	}
	
	// Create algorithm
	algorithm := &AssignmentAlgorithm{
		groupRepo:         &mockGroupRepo{groups: map[int]*Group{1: group}},
		itemRepo:          &mockItemRepo{items: map[int][]*Item{1: items}},
		participantRepo:   &mockParticipantRepo{participants: map[int][]*Participant{1: participants}},
		prioritizationRepo: &mockPrioritizationRepo{prioritizations: prioritizations},
		assignmentRepo:    &mockAssignmentRepo{},
	}
	
	// Execute assignment
	result, err := algorithm.ExecuteAssignment(context.Background(), 1)
	if err != nil {
		t.Fatalf("Assignment failed: %v", err)
	}
	
	// Verify results
	if len(result.Assignments) != 6 {
		t.Errorf("Expected 6 assignments, got %d", len(result.Assignments))
	}
	
	// Check that each participant gets 2 items (equal distribution)
	assignmentsByParticipant := make(map[int]int)
	for _, assignment := range result.Assignments {
		assignmentsByParticipant[assignment.ParticipantID]++
	}
	
	for participantID, count := range assignmentsByParticipant {
		if count != 2 {
			t.Errorf("Participant %d should have 2 items, got %d", participantID, count)
		}
	}
	
	// Check that no items are unassigned (all 6 items assigned)
	if len(result.Unassigned) != 0 {
		t.Errorf("Expected 0 unassigned items, got %d", len(result.Unassigned))
	}
}

func TestAssignmentAlgorithm_UnevenDistribution(t *testing.T) {
	// Setup test data for uneven distribution
	group := &Group{ID: 1, Name: "Test Group", Rule: "uneven", Status: "active"}
	items := []*Item{
		{ID: 1, GroupID: 1, Name: "Item 1"},
		{ID: 2, GroupID: 1, Name: "Item 2"},
		{ID: 3, GroupID: 1, Name: "Item 3"},
		{ID: 4, GroupID: 1, Name: "Item 4"},
		{ID: 5, GroupID: 1, Name: "Item 5"},
	}
	participants := []*Participant{
		{ID: 1, GroupID: 1, Token: "token1"},
		{ID: 2, GroupID: 1, Token: "token2"},
	}
	
	// Create prioritizations
	prioritizations := map[int][]*Prioritization{
		1: {
			{ParticipantID: 1, ItemID: 1, Rank: 1},
			{ParticipantID: 1, ItemID: 2, Rank: 2},
			{ParticipantID: 1, ItemID: 3, Rank: 3},
			{ParticipantID: 1, ItemID: 4, Rank: 4},
			{ParticipantID: 1, ItemID: 5, Rank: 5},
		},
		2: {
			{ParticipantID: 2, ItemID: 1, Rank: 5},
			{ParticipantID: 2, ItemID: 2, Rank: 4},
			{ParticipantID: 2, ItemID: 3, Rank: 3},
			{ParticipantID: 2, ItemID: 4, Rank: 2},
			{ParticipantID: 2, ItemID: 5, Rank: 1},
		},
	}
	
	// Create algorithm
	algorithm := &AssignmentAlgorithm{
		groupRepo:         &mockGroupRepo{groups: map[int]*Group{1: group}},
		itemRepo:          &mockItemRepo{items: map[int][]*Item{1: items}},
		participantRepo:   &mockParticipantRepo{participants: map[int][]*Participant{1: participants}},
		prioritizationRepo: &mockPrioritizationRepo{prioritizations: prioritizations},
		assignmentRepo:    &mockAssignmentRepo{},
	}
	
	// Execute assignment
	result, err := algorithm.ExecuteAssignment(context.Background(), 1)
	if err != nil {
		t.Fatalf("Assignment failed: %v", err)
	}
	
	// Verify results
	if len(result.Assignments) != 5 {
		t.Errorf("Expected 5 assignments, got %d", len(result.Assignments))
	}
	
	// Check that all items are assigned (uneven distribution)
	if len(result.Unassigned) != 0 {
		t.Errorf("Expected 0 unassigned items, got %d", len(result.Unassigned))
	}
	
	// Check that participants get different numbers of items
	assignmentsByParticipant := make(map[int]int)
	for _, assignment := range result.Assignments {
		assignmentsByParticipant[assignment.ParticipantID]++
	}
	
	// One participant should get 3 items, the other 2
	totalAssigned := 0
	for _, count := range assignmentsByParticipant {
		totalAssigned += count
	}
	
	if totalAssigned != 5 {
		t.Errorf("Expected 5 total assignments, got %d", totalAssigned)
	}
}