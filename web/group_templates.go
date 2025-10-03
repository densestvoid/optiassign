package web

import (
	"fmt"
	"optiassign/domain"
	"strconv"
	"strings"
)

// RenderGroupsList renders the groups list
func RenderGroupsList(groups []*domain.Group) string {
	if len(groups) == 0 {
		return `
			<div class="row">
				<div class="col-12">
					<div class="text-center">
						<h2>My Groups</h2>
						<p class="text-muted">You haven't created any groups yet.</p>
						<a href="/groups/new" class="btn btn-primary">Create Your First Group</a>
					</div>
				</div>
			</div>
		`
	}

	var groupsHTML strings.Builder
	groupsHTML.WriteString(`
		<div class="row">
			<div class="col-12">
				<div class="d-flex justify-content-between align-items-center mb-4">
					<h2>My Groups</h2>
					<a href="/groups/new" class="btn btn-primary">Create New Group</a>
				</div>
				<div class="row">
	`)

	for _, group := range groups {
		statusBadge := getStatusBadge(group.Status)
		groupsHTML.WriteString(fmt.Sprintf(`
					<div class="col-md-6 col-lg-4 mb-3">
						<div class="card">
							<div class="card-body">
								<h5 class="card-title">%s</h5>
								<p class="card-text">
									<span class="badge %s">%s</span>
									<span class="badge bg-secondary">%s</span>
								</p>
								<div class="d-grid gap-2">
									<a href="/groups/%d" class="btn btn-outline-primary">Manage Group</a>
								</div>
							</div>
						</div>
					</div>
		`, group.Name, statusBadge.Class, statusBadge.Text, strings.Title(group.Rule), group.ID))
	}

	groupsHTML.WriteString(`
				</div>
			</div>
		</div>
	`)

	return groupsHTML.String()
}

// RenderGroupDetails renders group details with items and participants
func RenderGroupDetails(group *domain.Group, items []*domain.Item, participants []*domain.Participant) string {
	statusBadge := getStatusBadge(group.Status)
	
	var content strings.Builder
	content.WriteString(fmt.Sprintf(`
		<div class="row">
			<div class="col-12">
				<div class="d-flex justify-content-between align-items-center mb-4">
					<div>
						<h2>%s</h2>
						<span class="badge %s">%s</span>
						<span class="badge bg-secondary">%s</span>
					</div>
					<div>
						<a href="/groups" class="btn btn-secondary">Back to Groups</a>
					</div>
				</div>
			</div>
		</div>
	`, group.Name, statusBadge.Class, statusBadge.Text, strings.Title(group.Rule)))

	// Items section
	content.WriteString(`
		<div class="row mb-4">
			<div class="col-12">
				<div class="card">
					<div class="card-header">
						<h5>Items (` + strconv.Itoa(len(items)) + `)</h5>
					</div>
					<div class="card-body">
	`)

	if len(items) == 0 {
		content.WriteString(`<p class="text-muted">No items added yet.</p>`)
	} else {
		content.WriteString(`<div class="list-group">`)
		for _, item := range items {
			assignedBadge := ""
			if item.IsAssigned {
				assignedBadge = `<span class="badge bg-success">Assigned</span>`
			}
			content.WriteString(fmt.Sprintf(`
				<div class="list-group-item d-flex justify-content-between align-items-center">
					<div>
						<h6 class="mb-1">%s</h6>
						%s
					</div>
					%s
				</div>
			`, item.Name, item.Description, assignedBadge))
		}
		content.WriteString(`</div>`)
	}

	content.WriteString(`
					</div>
				</div>
			</div>
		</div>
	`)

	// Participants section
	content.WriteString(`
		<div class="row">
			<div class="col-12">
				<div class="card">
					<div class="card-header">
						<h5>Participants (` + strconv.Itoa(len(participants)) + `)</h5>
					</div>
					<div class="card-body">
	`)

	if len(participants) == 0 {
		content.WriteString(`
			<p class="text-muted">No participants added yet.</p>
			<div class="mt-3">
				<button class="btn btn-primary" hx-get="/groups/` + strconv.Itoa(group.ID) + `/add-participant" hx-target="#add-participant-form" hx-swap="innerHTML">
					Add Participant
				</button>
			</div>
			<div id="add-participant-form"></div>
		`)
	} else {
		content.WriteString(`<div class="list-group">`)
		for _, participant := range participants {
			submittedBadge := ""
			if participant.HasSubmitted {
				submittedBadge = `<span class="badge bg-success">Submitted</span>`
			} else {
				submittedBadge = `<span class="badge bg-warning">Pending</span>`
			}
			content.WriteString(fmt.Sprintf(`
				<div class="list-group-item d-flex justify-content-between align-items-center">
					<div>
						<h6 class="mb-1">Participant #%d</h6>
						<small class="text-muted">Token: %s</small>
					</div>
					%s
				</div>
			`, participant.ID, participant.Token, submittedBadge))
		}
		content.WriteString(`</div>`)
		
		content.WriteString(`
			<div class="mt-3">
				<button class="btn btn-primary" hx-get="/groups/` + strconv.Itoa(group.ID) + `/add-participant" hx-target="#add-participant-form" hx-swap="innerHTML">
					Add Another Participant
				</button>
			</div>
			<div id="add-participant-form"></div>
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

// StatusBadge represents a status badge
type StatusBadge struct {
	Class string
	Text  string
}

// getStatusBadge returns the appropriate badge for a status
func getStatusBadge(status string) StatusBadge {
	switch status {
	case "draft":
		return StatusBadge{Class: "bg-secondary", Text: "Draft"}
	case "prioritizing":
		return StatusBadge{Class: "bg-warning", Text: "Prioritizing"}
	case "completed":
		return StatusBadge{Class: "bg-success", Text: "Completed"}
	default:
		return StatusBadge{Class: "bg-secondary", Text: strings.Title(status)}
	}
}