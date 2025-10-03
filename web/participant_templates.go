package web

import (
	"fmt"
	"optiassign/domain"
	"strconv"
	"strings"
)

// RenderParticipantView renders the participant access view
func RenderParticipantView(group *domain.Group, items []*domain.Item, participant *domain.Participant, hasSubmitted bool) string {
	var content strings.Builder
	
	content.WriteString(fmt.Sprintf(`
		<div class="row justify-content-center">
			<div class="col-md-8">
				<div class="card">
					<div class="card-header text-center">
						<h4>%s</h4>
						<p class="text-muted mb-0">Participant Access</p>
					</div>
					<div class="card-body">
	`, group.Name))

	if hasSubmitted {
		content.WriteString(`
			<div class="alert alert-success" role="alert">
				<h5 class="alert-heading">Priorities Submitted!</h5>
				<p>You have successfully submitted your priorities for this group.</p>
			</div>
		`)
	} else {
		content.WriteString(`
			<div class="alert alert-info" role="alert">
				<h5 class="alert-heading">Submit Your Priorities</h5>
				<p>Please rank all items in order of your preference (1 = most preferred).</p>
			</div>
		`)
	}

	// Items section
	content.WriteString(fmt.Sprintf(`
		<div class="mb-4">
			<h5>Items to Rank (%d total)</h5>
	`, len(items)))

	if len(items) == 0 {
		content.WriteString(`<p class="text-muted">No items available.</p>`)
	} else {
		content.WriteString(`<div class="list-group">`)
		for _, item := range items {
			content.WriteString(fmt.Sprintf(`
				<div class="list-group-item">
					<h6 class="mb-1">%s</h6>
					%s
				</div>
			`, item.Name, item.Description))
		}
		content.WriteString(`</div>`)
	}

	content.WriteString(`</div>`)

	// Action buttons
	if hasSubmitted {
		content.WriteString(`
			<div class="d-grid gap-2">
				<button class="btn btn-outline-primary" disabled>
					Priorities Already Submitted
				</button>
			</div>
		`)
	} else {
		content.WriteString(fmt.Sprintf(`
			<div class="d-grid gap-2">
				<button class="btn btn-primary" hx-get="/participant/%s/priorities" hx-target="#priority-form" hx-swap="innerHTML">
					Submit Priorities
				</button>
			</div>
			<div id="priority-form" class="mt-3"></div>
		`, participant.Token))
	}

	content.WriteString(`
					</div>
				</div>
			</div>
		</div>
	`)

	return content.String()
}

// RenderPriorityForm renders the priority submission form
func RenderPriorityForm(items []*domain.Item, token string) string {
	var content strings.Builder
	
	content.WriteString(`
		<div class="card">
			<div class="card-header">
				<h5>Submit Your Priorities</h5>
				<p class="text-muted mb-0">Rank all items from 1 (most preferred) to ` + strconv.Itoa(len(items)) + ` (least preferred)</p>
			</div>
			<div class="card-body">
				<form hx-post="/participant/` + token + `/priorities" hx-target="#priority-form" hx-swap="outerHTML">
	`)

	// Create form fields for each item
	for i, item := range items {
		content.WriteString(fmt.Sprintf(`
			<div class="mb-3">
				<label for="priority_%d" class="form-label">%d. %s</label>
				%s
				<select class="form-select" id="priority_%d" name="priority_%d" required>
					<option value="">Select rank...</option>
		`, item.ID, i+1, item.Name, item.Description, item.ID, item.ID))
		
		// Add rank options
		for rank := 1; rank <= len(items); rank++ {
			content.WriteString(fmt.Sprintf(`<option value="%d">%d</option>`, rank, rank))
		}
		
		content.WriteString(`
				</select>
			</div>
		`)
	}

	content.WriteString(`
					<div class="d-grid gap-2 d-md-flex justify-content-md-end">
						<button type="submit" class="btn btn-primary">Submit Priorities</button>
					</div>
				</form>
			</div>
		</div>
	`)

	return content.String()
}

// RenderAssignmentResults renders the assignment results
func RenderAssignmentResults(result *domain.AssignmentResult, group *domain.Group, participants []*domain.Participant, items []*domain.Item) string {
	var content strings.Builder
	
	content.WriteString(fmt.Sprintf(`
		<div class="row">
			<div class="col-12">
				<div class="alert alert-success" role="alert">
					<h4 class="alert-heading">Assignment Complete!</h4>
					<p>The assignment has been completed with %d rounds.</p>
				</div>
			</div>
		</div>
	`, result.Rounds))

	// Participant order
	content.WriteString(`
		<div class="row mb-4">
			<div class="col-12">
				<div class="card">
					<div class="card-header">
						<h5>Draft Order</h5>
					</div>
					<div class="card-body">
						<ol class="list-group list-group-numbered">
	`)

	for _, participantID := range result.ParticipantOrder {
		// Find participant name
		participantName := "Participant"
		for _, p := range participants {
			if p.ID == participantID {
				participantName = fmt.Sprintf("Participant #%d", p.ID)
				break
			}
		}
		
		content.WriteString(fmt.Sprintf(`
			<li class="list-group-item d-flex justify-content-between align-items-start">
				<div class="ms-2 me-auto">
					<div class="fw-bold">%s</div>
				</div>
			</li>
		`, participantName))
	}

	content.WriteString(`
						</ol>
					</div>
				</div>
			</div>
		</div>
	`)

	// Assignments by participant
	content.WriteString(`
		<div class="row mb-4">
			<div class="col-12">
				<div class="card">
					<div class="card-header">
						<h5>Assignments</h5>
					</div>
					<div class="card-body">
	`)

	// Group assignments by participant
	assignmentsByParticipant := make(map[int][]*domain.Assignment)
	for _, assignment := range result.Assignments {
		assignmentsByParticipant[assignment.ParticipantID] = append(assignmentsByParticipant[assignment.ParticipantID], assignment)
	}

	for _, participant := range participants {
		assignments := assignmentsByParticipant[participant.ID]
		content.WriteString(fmt.Sprintf(`
			<div class="mb-3">
				<h6>Participant #%d (%d items)</h6>
				<div class="list-group">
		`, participant.ID, len(assignments)))

		for _, assignment := range assignments {
			// Find item name
			itemName := "Unknown Item"
			for _, item := range items {
				if item.ID == assignment.ItemID {
					itemName = item.Name
					break
				}
			}
			
			content.WriteString(fmt.Sprintf(`
				<div class="list-group-item">%s</div>
			`, itemName))
		}

		content.WriteString(`
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

	// Unassigned items
	if len(result.Unassigned) > 0 {
		content.WriteString(`
			<div class="row">
				<div class="col-12">
					<div class="card">
						<div class="card-header">
							<h5>Unassigned Items</h5>
						</div>
						<div class="card-body">
							<div class="list-group">
		`)

		for _, item := range result.Unassigned {
			content.WriteString(fmt.Sprintf(`
				<div class="list-group-item">%s</div>
			`, item.Name))
		}

		content.WriteString(`
							</div>
						</div>
					</div>
				</div>
			</div>
		`)
	}

	return content.String()
}