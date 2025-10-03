package api

import (
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"optiassign/domain"
	"optiassign/web"
)

// GroupHandler handles group-related HTTP requests
type GroupHandler struct {
	groupService      *domain.GroupService
	userRepo          domain.UserRepository
	assignmentAlgorithm *domain.AssignmentAlgorithm
	emailService      *domain.EmailService
}

// NewGroupHandler creates a new group handler
func NewGroupHandler(groupService *domain.GroupService, userRepo domain.UserRepository, assignmentAlgorithm *domain.AssignmentAlgorithm, emailService *domain.EmailService) *GroupHandler {
	return &GroupHandler{
		groupService:        groupService,
		userRepo:            userRepo,
		assignmentAlgorithm: assignmentAlgorithm,
		emailService:        emailService,
	}
}

// ListGroups shows the user's groups
func (h *GroupHandler) ListGroups(w http.ResponseWriter, r *http.Request) {
	session := GetSessionFromContext(r)
	if session == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get user's groups
	groups, err := h.groupService.GetUserGroups(r.Context(), session.UserID)
	if err != nil {
		http.Error(w, "Failed to get groups", http.StatusInternalServerError)
		return
	}

	data := web.PageData{
		Title:   "My Groups - OptiAssign",
		Session: session,
		Content: template.HTML(web.RenderGroupsList(groups)),
	}

	web.RenderTemplate(w, "base.html", data)
}

// NewGroupForm shows the create group form
func (h *GroupHandler) NewGroupForm(w http.ResponseWriter, r *http.Request) {
	session := GetSessionFromContext(r)
	if session == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	content := `
		<div class="row justify-content-center">
			<div class="col-md-8">
				<div class="card">
					<div class="card-header">
						<h4>Create New Group</h4>
					</div>
					<div class="card-body">
						<form hx-post="/groups" hx-target="#group-form" hx-swap="outerHTML">
							<div class="mb-3">
								<label for="name" class="form-label">Group Name</label>
								<input type="text" class="form-control" id="name" name="name" required>
							</div>
							<div class="mb-3">
								<label for="rule" class="form-label">Distribution Rule</label>
								<select class="form-select" id="rule" name="rule" required>
									<option value="equal">Equal Distribution</option>
									<option value="exhaust">Exhaust All Items</option>
								</select>
							</div>
							<div class="mb-3">
								<label for="items" class="form-label">Items (one per line)</label>
								<textarea class="form-control" id="items" name="items" rows="5" required placeholder="Item 1&#10;Item 2&#10;Item 3"></textarea>
							</div>
							<div class="d-grid gap-2 d-md-flex justify-content-md-end">
								<a href="/dashboard" class="btn btn-secondary me-md-2">Cancel</a>
								<button type="submit" class="btn btn-primary">Create Group</button>
							</div>
						</form>
					</div>
				</div>
			</div>
		</div>
	`

	data := web.PageData{
		Title:   "Create Group - OptiAssign",
		Session: session,
		Content: template.HTML(content),
	}

	web.RenderTemplate(w, "base.html", data)
}

// CreateGroup creates a new group
func (h *GroupHandler) CreateGroup(w http.ResponseWriter, r *http.Request) {
	session := GetSessionFromContext(r)
	if session == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse form
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	name := strings.TrimSpace(r.FormValue("name"))
	rule := r.FormValue("rule")
	itemsText := strings.TrimSpace(r.FormValue("items"))

	// Validate inputs
	if name == "" {
		h.renderFormError(w, "Group name is required")
		return
	}
	if rule != "equal" && rule != "exhaust" {
		h.renderFormError(w, "Invalid distribution rule")
		return
	}
	if itemsText == "" {
		h.renderFormError(w, "At least one item is required")
		return
	}

	// Parse items
	itemLines := strings.Split(itemsText, "\n")
	var itemNames []string
	for _, line := range itemLines {
		itemName := strings.TrimSpace(line)
		if itemName != "" {
			itemNames = append(itemNames, itemName)
		}
	}

	if len(itemNames) == 0 {
		h.renderFormError(w, "At least one item is required")
		return
	}

	// Create group
	group, err := h.groupService.CreateGroupWithItems(r.Context(), session.UserID, name, rule, itemNames)
	if err != nil {
		h.renderFormError(w, "Failed to create group: "+err.Error())
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	
	successHTML := `
		<div class="alert alert-success" role="alert">
			<h4 class="alert-heading">Group Created Successfully!</h4>
			<p>Your group "<strong>` + group.Name + `</strong>" has been created with ` + strconv.Itoa(len(itemNames)) + ` items.</p>
			<hr>
			<div class="d-grid gap-2 d-md-flex justify-content-md-end">
				<a href="/groups/` + strconv.Itoa(group.ID) + `" class="btn btn-primary">Manage Group</a>
				<a href="/dashboard" class="btn btn-secondary">Back to Dashboard</a>
			</div>
		</div>
	`
	w.Write([]byte(successHTML))
}

// ViewGroup shows a group with its details
func (h *GroupHandler) ViewGroup(w http.ResponseWriter, r *http.Request) {
	session := GetSessionFromContext(r)
	if session == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	groupIDStr := chi.URLParam(r, "id")
	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		http.Error(w, "Invalid group ID", http.StatusBadRequest)
		return
	}

	// Get group details
	group, items, participants, err := h.groupService.GetGroupWithDetails(r.Context(), groupID)
	if err != nil {
		http.Error(w, "Failed to get group: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if user owns the group
	if group.OwnerUserID != session.UserID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	content := web.RenderGroupDetails(group, items, participants)
	
	data := web.PageData{
		Title:   "Group: " + group.Name + " - OptiAssign",
		Session: session,
		Content: template.HTML(content),
	}

	web.RenderTemplate(w, "base.html", data)
}

// ExecuteAssignment runs the assignment algorithm
func (h *GroupHandler) ExecuteAssignment(w http.ResponseWriter, r *http.Request) {
	session := GetSessionFromContext(r)
	if session == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	groupIDStr := chi.URLParam(r, "id")
	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		http.Error(w, "Invalid group ID", http.StatusBadRequest)
		return
	}

	// Get group to verify ownership
	group, err := h.groupService.GetGroupByID(r.Context(), groupID)
	if err != nil {
		http.Error(w, "Failed to get group: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if group.OwnerUserID != session.UserID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Execute assignment
	result, err := h.assignmentAlgorithm.ExecuteAssignment(r.Context(), groupID)
	if err != nil {
		h.renderFormError(w, "Failed to execute assignment: "+err.Error())
		return
	}

	// Send email notifications to participants
	participants, err := h.groupService.GetGroupParticipants(r.Context(), groupID)
	if err == nil {
		baseURL := "http://localhost:8080" // In production, this should come from config
		for _, participant := range participants {
			h.emailService.SendAssignmentNotification(participant, group, result.Assignments, baseURL)
		}
	}

	// Return success response
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	
	successHTML := `
		<div class="alert alert-success" role="alert">
			<h4 class="alert-heading">Assignment Complete!</h4>
			<p>The assignment has been executed successfully. All participants have been notified.</p>
			<div class="mt-3">
				<a href="/groups/` + strconv.Itoa(groupID) + `" class="btn btn-primary">View Results</a>
			</div>
		</div>
	`
	w.Write([]byte(successHTML))
}

// renderFormError renders a form error response
func (h *GroupHandler) renderFormError(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusBadRequest)
	
	errorHTML := `
		<div class="alert alert-danger" role="alert">
			<strong>Error:</strong> ` + message + `
		</div>
	`
	w.Write([]byte(errorHTML))
}