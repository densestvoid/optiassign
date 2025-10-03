package api

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"optiassign/domain"
	"optiassign/web"
)

// ParticipantHandler handles participant-related HTTP requests
type ParticipantHandler struct {
	groupService        *domain.GroupService
	prioritizationRepo  domain.PrioritizationRepository
	participantRepo     domain.ParticipantRepository
	itemRepo            domain.ItemRepository
	assignmentService   *domain.AssignmentService
}

// NewParticipantHandler creates a new participant handler
func NewParticipantHandler(
	groupService *domain.GroupService,
	prioritizationRepo domain.PrioritizationRepository,
	participantRepo domain.ParticipantRepository,
	itemRepo domain.ItemRepository,
	assignmentService *domain.AssignmentService,
) *ParticipantHandler {
	return &ParticipantHandler{
		groupService:        groupService,
		prioritizationRepo:  prioritizationRepo,
		participantRepo:     participantRepo,
		itemRepo:            itemRepo,
		assignmentService:   assignmentService,
	}
}

// ParticipantAccess handles tokenized participant access
func (h *ParticipantHandler) ParticipantAccess(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	if token == "" {
		http.Error(w, "Token required", http.StatusBadRequest)
		return
	}

	// Get participant by token
	participant, err := h.participantRepo.GetByToken(token)
	if err != nil {
		http.Error(w, "Failed to get participant", http.StatusInternalServerError)
		return
	}
	if participant == nil {
		http.Error(w, "Invalid token", http.StatusNotFound)
		return
	}

	// Get group details
	group, items, _, err := h.groupService.GetGroupWithDetails(r.Context(), participant.GroupID)
	if err != nil {
		http.Error(w, "Failed to get group details", http.StatusInternalServerError)
		return
	}

	// Check if participant has already submitted
	hasSubmitted, err := h.prioritizationRepo.HasSubmitted(participant.ID)
	if err != nil {
		http.Error(w, "Failed to check submission status", http.StatusInternalServerError)
		return
	}

	content := web.RenderParticipantView(group, items, participant, hasSubmitted)
	
	data := web.PageData{
		Title:   "Participant Access - " + group.Name,
		Content: template.HTML(content),
	}

	web.RenderTemplate(w, "base.html", data)
}

// PriorityForm shows the priority submission form
func (h *ParticipantHandler) PriorityForm(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	if token == "" {
		http.Error(w, "Token required", http.StatusBadRequest)
		return
	}

	// Get participant by token
	participant, err := h.participantRepo.GetByToken(token)
	if err != nil {
		http.Error(w, "Failed to get participant", http.StatusInternalServerError)
		return
	}
	if participant == nil {
		http.Error(w, "Invalid token", http.StatusNotFound)
		return
	}

	// Get items
	items, err := h.itemRepo.GetByGroupID(participant.GroupID)
	if err != nil {
		http.Error(w, "Failed to get items", http.StatusInternalServerError)
		return
	}

	// Check if already submitted
	hasSubmitted, err := h.prioritizationRepo.HasSubmitted(participant.ID)
	if err != nil {
		http.Error(w, "Failed to check submission status", http.StatusInternalServerError)
		return
	}

	if hasSubmitted {
		http.Error(w, "Priorities already submitted", http.StatusBadRequest)
		return
	}

	content := web.RenderPriorityForm(items, token)
	
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(content))
}

// SubmitPriorities handles priority submission
func (h *ParticipantHandler) SubmitPriorities(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	token := chi.URLParam(r, "token")
	if token == "" {
		http.Error(w, "Token required", http.StatusBadRequest)
		return
	}

	// Get participant by token
	participant, err := h.participantRepo.GetByToken(token)
	if err != nil {
		http.Error(w, "Failed to get participant", http.StatusInternalServerError)
		return
	}
	if participant == nil {
		http.Error(w, "Invalid token", http.StatusNotFound)
		return
	}

	// Check if already submitted
	hasSubmitted, err := h.prioritizationRepo.HasSubmitted(participant.ID)
	if err != nil {
		http.Error(w, "Failed to check submission status", http.StatusInternalServerError)
		return
	}

	if hasSubmitted {
		http.Error(w, "Priorities already submitted", http.StatusBadRequest)
		return
	}

	// Parse form
	if err := r.ParseForm(); err != nil {
		h.renderPriorityError(w, "Invalid form data")
		return
	}

	// Get items
	items, err := h.itemRepo.GetByGroupID(participant.GroupID)
	if err != nil {
		h.renderPriorityError(w, "Failed to get items")
		return
	}

	// Parse priorities
	priorities := make(map[int]int) // itemID -> rank
	for _, item := range items {
		rankStr := r.FormValue("priority_" + strconv.Itoa(item.ID))
		if rankStr == "" {
			h.renderPriorityError(w, "All items must be ranked")
			return
		}
		
		rank, err := strconv.Atoi(rankStr)
		if err != nil {
			h.renderPriorityError(w, "Invalid rank for item: " + item.Name)
			return
		}
		
		priorities[item.ID] = rank
	}

	// Validate priorities (must be unique ranks from 1 to len(items))
	usedRanks := make(map[int]bool)
	for _, rank := range priorities {
		if rank < 1 || rank > len(items) {
			h.renderPriorityError(w, "Ranks must be between 1 and " + strconv.Itoa(len(items)))
			return
		}
		if usedRanks[rank] {
			h.renderPriorityError(w, "Each rank can only be used once")
			return
		}
		usedRanks[rank] = true
	}

	// Clear existing priorities
	if err := h.prioritizationRepo.DeleteByParticipantID(participant.ID); err != nil {
		h.renderPriorityError(w, "Failed to clear existing priorities")
		return
	}

	// Save new priorities
	for itemID, rank := range priorities {
		prioritization := &domain.Prioritization{
			ParticipantID: participant.ID,
			ItemID:        itemID,
			Rank:          rank,
		}
		
		if err := h.prioritizationRepo.Create(prioritization); err != nil {
			h.renderPriorityError(w, "Failed to save priorities")
			return
		}
	}

	// Update participant status
	participant.HasSubmitted = true
	if err := h.participantRepo.Update(participant); err != nil {
		h.renderPriorityError(w, "Failed to update participant status")
		return
	}

	// Check if all participants have submitted and trigger assignment if ready
	if err := h.assignmentService.CheckAndTriggerAssignment(r.Context(), participant.GroupID); err != nil {
		// Log error but don't fail the request - assignment will be retried
		log.Printf("Failed to check/trigger assignment for group %d: %v", participant.GroupID, err)
	}

	// Return success response
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	
	successHTML := `
		<div class="alert alert-success" role="alert">
			<h4 class="alert-heading">Priorities Submitted Successfully!</h4>
			<p>Your priorities have been saved. You can view them below.</p>
		</div>
		<div class="mt-3">
			<a href="/participant/` + token + `" class="btn btn-primary">View Your Priorities</a>
		</div>
	`
	w.Write([]byte(successHTML))
}

// AssignmentStatus returns the assignment status for a group
func (h *ParticipantHandler) AssignmentStatus(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	if token == "" {
		http.Error(w, "Token required", http.StatusBadRequest)
		return
	}

	// Get participant by token
	participant, err := h.participantRepo.GetByToken(token)
	if err != nil {
		http.Error(w, "Failed to get participant", http.StatusInternalServerError)
		return
	}
	if participant == nil {
		http.Error(w, "Invalid token", http.StatusNotFound)
		return
	}

	// Get assignment status
	status, err := h.assignmentService.GetAssignmentStatus(r.Context(), participant.GroupID)
	if err != nil {
		http.Error(w, "Failed to get assignment status", http.StatusInternalServerError)
		return
	}

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	// Simple JSON response (in a real app, use json.Marshal)
	jsonResponse := fmt.Sprintf(`{
		"group_id": %d,
		"group_status": "%s",
		"total_participants": %d,
		"submitted_count": %d,
		"all_submitted": %t,
		"is_completed": %t
	}`, status.GroupID, status.GroupStatus, status.TotalParticipants, status.SubmittedCount, status.AllSubmitted, status.IsCompleted)
	
	w.Write([]byte(jsonResponse))
}

// renderPriorityError renders a priority form error response
func (h *ParticipantHandler) renderPriorityError(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusBadRequest)
	
	errorHTML := `
		<div class="alert alert-danger" role="alert">
			<strong>Error:</strong> ` + message + `
		</div>
	`
	w.Write([]byte(errorHTML))
}