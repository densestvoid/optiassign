package api

import (
	"encoding/json"
	"net/http"
	"optiassign/domain"
	"optiassign/web"
)

// MonitoringHandler handles monitoring-related HTTP requests
type MonitoringHandler struct {
	metricsCollector *domain.MetricsCollector
}

// NewMonitoringHandler creates a new monitoring handler
func NewMonitoringHandler(metricsCollector *domain.MetricsCollector) *MonitoringHandler {
	return &MonitoringHandler{
		metricsCollector: metricsCollector,
	}
}

// Dashboard renders the monitoring dashboard
func (h *MonitoringHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	// Get current metrics
	metrics := h.metricsCollector.GetMetrics()
	uptime := h.metricsCollector.GetUptime()
	requestsPerMinute := h.metricsCollector.GetRequestsPerMinute()

	// Render dashboard
	if err := web.RenderMonitoringDashboard(w, metrics, uptime, requestsPerMinute); err != nil {
		http.Error(w, "Failed to render monitoring dashboard", http.StatusInternalServerError)
		return
	}
}

// MetricsAPI returns metrics as JSON
func (h *MonitoringHandler) MetricsAPI(w http.ResponseWriter, r *http.Request) {
	metrics := h.metricsCollector.GetMetrics()
	uptime := h.metricsCollector.GetUptime()
	requestsPerMinute := h.metricsCollector.GetRequestsPerMinute()

	response := map[string]interface{}{
		"metrics":           metrics,
		"uptime_seconds":    uptime.Seconds(),
		"requests_per_minute": requestsPerMinute,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	json.NewEncoder(w).Encode(response)
}