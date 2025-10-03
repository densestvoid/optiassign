package web

import (
	"fmt"
	"optiassign/domain"
	"strings"
)

// RenderFullAssignmentResults renders the full assignment results for group owners
func RenderFullAssignmentResults(group *domain.Group, assignments []*domain.Assignment, participants []*domain.Participant, items []*domain.Item, status *domain.AssignmentStatus) string {
	var content strings.Builder
	
	content.WriteString(fmt.Sprintf(`
		<div class="row">
			<div class="col-12">
				<div class="card">
					<div class="card-header">
						<h4>Assignment Results: %s</h4>
						<p class="text-muted mb-0">Distribution Rule: %s</p>
					</div>
					<div class="card-body">
	`, group.Name, group.Rule))

	// Status summary
	content.WriteString(fmt.Sprintf(`
		<div class="row mb-4">
			<div class="col-md-3">
				<div class="card text-center">
					<div class="card-body">
						<h5 class="card-title">%d</h5>
						<p class="card-text">Total Participants</p>
					</div>
				</div>
			</div>
			<div class="col-md-3">
				<div class="card text-center">
					<div class="card-body">
						<h5 class="card-title">%d</h5>
						<p class="card-text">Items Assigned</p>
					</div>
				</div>
			</div>
			<div class="col-md-3">
				<div class="card text-center">
					<div class="card-body">
						<h5 class="card-title">%d</h5>
						<p class="card-text">Items Remaining</p>
					</div>
				</div>
			</div>
			<div class="col-md-3">
				<div class="card text-center">
					<div class="card-body">
						<h5 class="card-title">%s</h5>
						<p class="card-text">Status</p>
					</div>
				</div>
			</div>
		</div>
	`, status.TotalParticipants, len(assignments), len(items)-len(assignments), status.GroupStatus))

	// Assignments by participant
	content.WriteString(`
		<div class="row">
			<div class="col-12">
				<h5>Assignments by Participant</h5>
	`)

	// Group assignments by participant
	assignmentsByParticipant := make(map[int][]*domain.Assignment)
	for _, assignment := range assignments {
		assignmentsByParticipant[assignment.ParticipantID] = append(assignmentsByParticipant[assignment.ParticipantID], assignment)
	}

	for _, participant := range participants {
		participantAssignments := assignmentsByParticipant[participant.ID]
		content.WriteString(fmt.Sprintf(`
			<div class="card mb-3">
				<div class="card-header">
					<h6 class="mb-0">Participant #%d (%d items)</h6>
				</div>
				<div class="card-body">
		`, participant.ID, len(participantAssignments)))

		if len(participantAssignments) == 0 {
			content.WriteString(`<p class="text-muted">No items assigned</p>`)
		} else {
			content.WriteString(`<div class="list-group">`)
			for _, assignment := range participantAssignments {
				// Find item details
				itemName := "Unknown Item"
				itemDescription := ""
				for _, item := range items {
					if item.ID == assignment.ItemID {
						itemName = item.Name
						itemDescription = item.Description
						break
					}
				}
				
				content.WriteString(fmt.Sprintf(`
					<div class="list-group-item">
						<h6 class="mb-1">%s</h6>
						%s
					</div>
				`, itemName, itemDescription))
			}
			content.WriteString(`</div>`)
		}

		content.WriteString(`
				</div>
			</div>
		`)
	}

	content.WriteString(`
			</div>
		</div>
	`)

	// Unassigned items
	assignedItemIDs := make(map[int]bool)
	for _, assignment := range assignments {
		assignedItemIDs[assignment.ItemID] = true
	}

	var unassignedItems []*domain.Item
	for _, item := range items {
		if !assignedItemIDs[item.ID] {
			unassignedItems = append(unassignedItems, item)
		}
	}

	if len(unassignedItems) > 0 {
		content.WriteString(`
			<div class="row mt-4">
				<div class="col-12">
					<div class="card">
						<div class="card-header">
							<h5>Unassigned Items (%d)</h5>
						</div>
						<div class="card-body">
							<div class="list-group">
		`)

		for _, item := range unassignedItems {
			content.WriteString(fmt.Sprintf(`
				<div class="list-group-item">
					<h6 class="mb-1">%s</h6>
					%s
				</div>
			`, item.Name, item.Description))
		}

		content.WriteString(`
							</div>
						</div>
					</div>
				</div>
			</div>
		`)
	}

	content.WriteString(`
					</div>
				</div>
			</div>
		</div>
	`)

	return content.String()
}

// RenderParticipantResults renders assignment results for a participant
func RenderParticipantResults(group *domain.Group, assignments []*domain.Assignment, items []*domain.Item, status *domain.AssignmentStatus) string {
	var content strings.Builder
	
	content.WriteString(fmt.Sprintf(`
		<div class="row justify-content-center">
			<div class="col-md-8">
				<div class="card">
					<div class="card-header text-center">
						<h4>Your Assignment Results</h4>
						<p class="text-muted mb-0">%s</p>
					</div>
					<div class="card-body">
	`, group.Name))

	if status.IsCompleted {
		content.WriteString(`
			<div class="alert alert-success" role="alert">
				<h5 class="alert-heading">Assignment Complete!</h5>
				<p>You have been assigned the following items:</p>
			</div>
		`)
	} else {
		content.WriteString(`
			<div class="alert alert-info" role="alert">
				<h5 class="alert-heading">Assignment in Progress</h5>
				<p>%d of %d participants have submitted their priorities.</p>
			</div>
		`)
	}

	// Show assigned items
	if len(assignments) > 0 {
		content.WriteString(fmt.Sprintf(`
			<div class="mb-4">
				<h5>Your Assigned Items (%d)</h5>
				<div class="list-group">
		`, len(assignments)))

		for i, assignment := range assignments {
			// Find item details
			itemName := "Unknown Item"
			itemDescription := ""
			for _, item := range items {
				if item.ID == assignment.ItemID {
					itemName = item.Name
					itemDescription = item.Description
					break
				}
			}
			
			content.WriteString(fmt.Sprintf(`
				<div class="list-group-item">
					<div class="d-flex w-100 justify-content-between">
						<h6 class="mb-1">%d. %s</h6>
						<small>Assigned</small>
					</div>
					%s
				</div>
			`, i+1, itemName, itemDescription))
		}

		content.WriteString(`
				</div>
			</div>
		`)
	} else {
		content.WriteString(`
			<div class="alert alert-warning" role="alert">
				<p>No items have been assigned to you yet.</p>
			</div>
		`)
	}

	// Status information
	content.WriteString(fmt.Sprintf(`
		<div class="row">
			<div class="col-md-6">
				<div class="card">
					<div class="card-body text-center">
						<h5 class="card-title">%d</h5>
						<p class="card-text">Total Participants</p>
					</div>
				</div>
			</div>
			<div class="col-md-6">
				<div class="card">
					<div class="card-body text-center">
						<h5 class="card-title">%d</h5>
						<p class="card-text">Items Assigned to You</p>
					</div>
				</div>
			</div>
		</div>
	`, status.TotalParticipants, len(assignments)))

	content.WriteString(`
					</div>
				</div>
			</div>
		</div>
	`)

	return content.String()
}