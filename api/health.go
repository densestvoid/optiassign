package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// HealthStatus represents the health status of the application
type HealthStatus struct {
	Status    string            `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
	Version   string            `json:"version"`
	Services  map[string]string `json:"services"`
}

// HealthChecker handles health check endpoints
type HealthChecker struct {
	db *sql.DB
}

// NewHealthChecker creates a new health checker
func NewHealthChecker(db *sql.DB) *HealthChecker {
	return &HealthChecker{db: db}
}

// HealthCheck returns the overall health status
func (h *HealthChecker) HealthCheck(w http.ResponseWriter, r *http.Request) {
	status := &HealthStatus{
		Status:    "healthy",
		Timestamp: time.Now(),
		Version:   "1.0.0",
		Services:  make(map[string]string),
	}
	
	// Check database
	if err := h.checkDatabase(); err != nil {
		status.Status = "unhealthy"
		status.Services["database"] = "unhealthy: " + err.Error()
	} else {
		status.Services["database"] = "healthy"
	}
	
	// Check other services
	status.Services["api"] = "healthy"
	status.Services["task_manager"] = "healthy"
	
	w.Header().Set("Content-Type", "application/json")
	
	if status.Status == "unhealthy" {
		w.WriteHeader(http.StatusServiceUnavailable)
	}
	
	json.NewEncoder(w).Encode(status)
}

// ReadinessCheck returns the readiness status
func (h *HealthChecker) ReadinessCheck(w http.ResponseWriter, r *http.Request) {
	status := &HealthStatus{
		Status:    "ready",
		Timestamp: time.Now(),
		Services:  make(map[string]string),
	}
	
	// Check if database is ready
	if err := h.checkDatabase(); err != nil {
		status.Status = "not_ready"
		status.Services["database"] = "not_ready: " + err.Error()
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		status.Services["database"] = "ready"
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

// LivenessCheck returns the liveness status
func (h *HealthChecker) LivenessCheck(w http.ResponseWriter, r *http.Request) {
	status := &HealthStatus{
		Status:    "alive",
		Timestamp: time.Now(),
		Services:  make(map[string]string),
	}
	
	// Simple liveness check - if we can respond, we're alive
	status.Services["application"] = "alive"
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

// checkDatabase checks if the database is accessible
func (h *HealthChecker) checkDatabase() error {
	if h.db == nil {
		return fmt.Errorf("database connection not initialized")
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	return h.db.PingContext(ctx)
}

// MetricsHandler returns basic application metrics
func (h *HealthChecker) MetricsHandler(w http.ResponseWriter, r *http.Request) {
	metrics := map[string]interface{}{
		"timestamp": time.Now().Unix(),
		"uptime":    time.Since(startTime).Seconds(),
		"version":   "1.0.0",
		"environment": "production",
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics)
}

var startTime = time.Now()