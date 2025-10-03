package api

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"optiassign/domain"
	"optiassign/web"
)

// AssignmentHandler handles assignment results display
type AssignmentHandler struct {
	groupService      *domain.GroupService
	assignmentRepo    domain.AssignmentRepository
	participantRepo   domain.ParticipantRepository
	itemRepo          domain.ItemRepository
	assignmentService *domain.AssignmentService
}

// NewAssignmentHandler creates a new assignment handler
func NewAssignmentHandler(
	groupService *domain.GroupService,
	assignmentRepo domain.AssignmentRepository,
	participantRepo domain.ParticipantRepository,
	itemRepo domain.ItemRepository,
	assignmentService *domain.AssignmentService,
) *AssignmentHandler {
	return &AssignmentHandler{
		groupService:      groupService,
		assignmentRepo:    assignmentRepo,
		participantRepo:   participantRepo,
		itemRepo:          itemRepo,
		assignmentService: assignmentService,
	}
}

// ViewResults displays assignment results for a group
func (h *AssignmentHandler) ViewResults(w http.ResponseWriter, r *http.Request) {
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

	// Get assignment results
	assignments, err := h.assignmentRepo.GetByGroupID(groupID)
	if err != nil {
		http.Error(w, "Failed to get assignments: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Get participants and items
	participants, err := h.participantRepo.GetByGroupID(groupID)
	if err != nil {
		http.Error(w, "Failed to get participants: "+err.Error(), http.StatusInternalServerError)
		return
	}

	items, err := h.itemRepo.GetByGroupID(groupID)
	if err != nil {
		http.Error(w, "Failed to get items: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Get assignment status
	status, err := h.assignmentService.GetAssignmentStatus(r.Context(), groupID)
	if err != nil {
		http.Error(w, "Failed to get assignment status: "+err.Error(), http.StatusInternalServerError)
		return
	}

	content := web.RenderFullAssignmentResults(group, assignments, participants, items, status)
	
	data := web.PageData{
		Title:   "Assignment Results - " + group.Name,
		Content: template.HTML(content),
	}

	web.RenderTemplate(w, "base.html", data)
}

// ParticipantResults displays results for a participant
func (h *AssignmentHandler) ParticipantResults(w http.ResponseWriter, r *http.Request) {
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
	group, err := h.groupService.GetGroupByID(r.Context(), participant.GroupID)
	if err != nil {
		http.Error(w, "Failed to get group", http.StatusInternalServerError)
		return
	}

	// Get participant's assignments
	assignments, err := h.assignmentRepo.GetByParticipantID(participant.ID)
	if err != nil {
		http.Error(w, "Failed to get assignments", http.StatusInternalServerError)
		return
	}

	// Get items
	items, err := h.itemRepo.GetByGroupID(participant.GroupID)
	if err != nil {
		http.Error(w, "Failed to get items", http.StatusInternalServerError)
		return
	}

	// Get assignment status
	status, err := h.assignmentService.GetAssignmentStatus(r.Context(), participant.GroupID)
	if err != nil {
		http.Error(w, "Failed to get assignment status", http.StatusInternalServerError)
		return
	}

	content := web.RenderParticipantResults(group, assignments, items, status)
	
	data := web.PageData{
		Title:   "Your Assignment Results - " + group.Name,
		Content: template.HTML(content),
	}

	web.RenderTemplate(w, "base.html", data)
}