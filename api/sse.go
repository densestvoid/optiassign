package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// SSEEvent represents a Server-Sent Event
type SSEEvent struct {
	ID    string      `json:"id"`
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}

// SSEClient represents a connected SSE client
type SSEClient struct {
	ID       string
	GroupID  int
	Token    string
	Messages chan SSEEvent
	Context  context.Context
	Cancel   context.CancelFunc
}

// SSEManager manages Server-Sent Events
type SSEManager struct {
	clients map[string]*SSEClient
	mutex   sync.RWMutex
}

// NewSSEManager creates a new SSE manager
func NewSSEManager() *SSEManager {
	return &SSEManager{
		clients: make(map[string]*SSEClient),
	}
}

// AddClient adds a new SSE client
func (m *SSEManager) AddClient(clientID string, groupID int, token string) *SSEClient {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	ctx, cancel := context.WithCancel(context.Background())
	client := &SSEClient{
		ID:       clientID,
		GroupID:  groupID,
		Token:    token,
		Messages: make(chan SSEEvent, 10),
		Context:  ctx,
		Cancel:   cancel,
	}
	
	m.clients[clientID] = client
	return client
}

// RemoveClient removes an SSE client
func (m *SSEManager) RemoveClient(clientID string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	if client, exists := m.clients[clientID]; exists {
		client.Cancel()
		close(client.Messages)
		delete(m.clients, clientID)
	}
}

// BroadcastToGroup broadcasts an event to all clients in a group
func (m *SSEManager) BroadcastToGroup(groupID int, event SSEEvent) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	
	for _, client := range m.clients {
		if client.GroupID == groupID {
			select {
			case client.Messages <- event:
			default:
				// Channel full, skip this client
			}
		}
	}
}

// BroadcastToParticipant broadcasts an event to a specific participant
func (m *SSEManager) BroadcastToParticipant(token string, event SSEEvent) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	
	for _, client := range m.clients {
		if client.Token == token {
			select {
			case client.Messages <- event:
			default:
				// Channel full, skip this client
			}
		}
	}
}

// SSEHandler handles Server-Sent Events connections
func (m *SSEManager) SSEHandler(w http.ResponseWriter, r *http.Request) {
	// Get parameters
	groupIDStr := r.URL.Query().Get("group_id")
	token := r.URL.Query().Get("token")
	clientID := r.URL.Query().Get("client_id")
	
	if groupIDStr == "" || clientID == "" {
		http.Error(w, "Missing required parameters", http.StatusBadRequest)
		return
	}
	
	// Parse group ID
	var groupID int
	if _, err := fmt.Sscanf(groupIDStr, "%d", &groupID); err != nil {
		http.Error(w, "Invalid group ID", http.StatusBadRequest)
		return
	}
	
	// Set SSE headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Cache-Control")
	
	// Add client
	client := m.AddClient(clientID, groupID, token)
	defer m.RemoveClient(clientID)
	
	// Send initial connection event
	event := SSEEvent{
		ID:    fmt.Sprintf("%d", time.Now().Unix()),
		Event: "connected",
		Data:  map[string]interface{}{"message": "Connected to real-time updates"},
	}
	
	if err := m.writeSSEEvent(w, event); err != nil {
		return
	}
	
	// Keep connection alive and send events
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	
	for {
		select {
		case event := <-client.Messages:
			if err := m.writeSSEEvent(w, event); err != nil {
				return
			}
		case <-ticker.C:
			// Send keepalive
			keepalive := SSEEvent{
				ID:    fmt.Sprintf("%d", time.Now().Unix()),
				Event: "keepalive",
				Data:  map[string]interface{}{"timestamp": time.Now().Unix()},
			}
			if err := m.writeSSEEvent(w, keepalive); err != nil {
				return
			}
		case <-client.Context.Done():
			return
		case <-r.Context().Done():
			return
		}
	}
}

// writeSSEEvent writes an SSE event to the response
func (m *SSEManager) writeSSEEvent(w http.ResponseWriter, event SSEEvent) error {
	// Write SSE format
	fmt.Fprintf(w, "id: %s\n", event.ID)
	fmt.Fprintf(w, "event: %s\n", event.Event)
	
	// Convert data to JSON
	data, err := json.Marshal(event.Data)
	if err != nil {
		return err
	}
	
	fmt.Fprintf(w, "data: %s\n\n", string(data))
	
	// Flush the response
	if flusher, ok := w.(http.Flusher); ok {
		flusher.Flush()
	}
	
	return nil
}

// Global SSE manager instance
var sseManager = NewSSEManager()

// GetSSEManager returns the global SSE manager
func GetSSEManager() *SSEManager {
	return sseManager
}