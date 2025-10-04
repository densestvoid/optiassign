package domain

import (
	"sync"
	"time"
)

// Metrics represents application metrics
type Metrics struct {
	RequestsTotal     int64     `json:"requests_total"`
	RequestsByStatus  map[int]int64 `json:"requests_by_status"`
	ActiveUsers       int64     `json:"active_users"`
	GroupsCreated     int64     `json:"groups_created"`
	AssignmentsRun    int64     `json:"assignments_run"`
	EmailsSent        int64     `json:"emails_sent"`
	LastActivity      time.Time `json:"last_activity"`
	StartTime         time.Time `json:"start_time"`
	mutex             sync.RWMutex
}

// MetricsCollector collects and manages application metrics
type MetricsCollector struct {
	metrics *Metrics
}

// NewMetricsCollector creates a new metrics collector
func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{
		metrics: &Metrics{
			RequestsByStatus: make(map[int]int64),
			StartTime:        time.Now(),
		},
	}
}

// IncrementRequests increments the total request counter
func (mc *MetricsCollector) IncrementRequests() {
	mc.metrics.mutex.Lock()
	defer mc.metrics.mutex.Unlock()
	mc.metrics.RequestsTotal++
	mc.metrics.LastActivity = time.Now()
}

// IncrementRequestsByStatus increments the request counter for a specific status
func (mc *MetricsCollector) IncrementRequestsByStatus(status int) {
	mc.metrics.mutex.Lock()
	defer mc.metrics.mutex.Unlock()
	mc.metrics.RequestsByStatus[status]++
	mc.metrics.LastActivity = time.Now()
}

// IncrementActiveUsers increments the active users counter
func (mc *MetricsCollector) IncrementActiveUsers() {
	mc.metrics.mutex.Lock()
	defer mc.metrics.mutex.Unlock()
	mc.metrics.ActiveUsers++
	mc.metrics.LastActivity = time.Now()
}

// DecrementActiveUsers decrements the active users counter
func (mc *MetricsCollector) DecrementActiveUsers() {
	mc.metrics.mutex.Lock()
	defer mc.metrics.mutex.Unlock()
	if mc.metrics.ActiveUsers > 0 {
		mc.metrics.ActiveUsers--
	}
}

// IncrementGroupsCreated increments the groups created counter
func (mc *MetricsCollector) IncrementGroupsCreated() {
	mc.metrics.mutex.Lock()
	defer mc.metrics.mutex.Unlock()
	mc.metrics.GroupsCreated++
	mc.metrics.LastActivity = time.Now()
}

// IncrementAssignmentsRun increments the assignments run counter
func (mc *MetricsCollector) IncrementAssignmentsRun() {
	mc.metrics.mutex.Lock()
	defer mc.metrics.mutex.Unlock()
	mc.metrics.AssignmentsRun++
	mc.metrics.LastActivity = time.Now()
}

// IncrementEmailsSent increments the emails sent counter
func (mc *MetricsCollector) IncrementEmailsSent() {
	mc.metrics.mutex.Lock()
	defer mc.metrics.mutex.Unlock()
	mc.metrics.EmailsSent++
	mc.metrics.LastActivity = time.Now()
}

// GetMetrics returns the current metrics
func (mc *MetricsCollector) GetMetrics() *Metrics {
	mc.metrics.mutex.RLock()
	defer mc.metrics.mutex.RUnlock()
	
	// Return a copy to avoid race conditions
	metrics := &Metrics{
		RequestsTotal:    mc.metrics.RequestsTotal,
		RequestsByStatus: make(map[int]int64),
		ActiveUsers:      mc.metrics.ActiveUsers,
		GroupsCreated:    mc.metrics.GroupsCreated,
		AssignmentsRun:   mc.metrics.AssignmentsRun,
		EmailsSent:       mc.metrics.EmailsSent,
		LastActivity:     mc.metrics.LastActivity,
		StartTime:        mc.metrics.StartTime,
	}
	
	// Copy the requests by status map
	for status, count := range mc.metrics.RequestsByStatus {
		metrics.RequestsByStatus[status] = count
	}
	
	return metrics
}

// GetUptime returns the application uptime
func (mc *MetricsCollector) GetUptime() time.Duration {
	mc.metrics.mutex.RLock()
	defer mc.metrics.mutex.RUnlock()
	return time.Since(mc.metrics.StartTime)
}

// GetRequestsPerMinute calculates requests per minute
func (mc *MetricsCollector) GetRequestsPerMinute() float64 {
	mc.metrics.mutex.RLock()
	defer mc.metrics.mutex.RUnlock()
	
	uptime := time.Since(mc.metrics.StartTime)
	if uptime.Minutes() == 0 {
		return 0
	}
	
	return float64(mc.metrics.RequestsTotal) / uptime.Minutes()
}

// Global metrics collector instance
var globalMetrics = NewMetricsCollector()

// GetGlobalMetrics returns the global metrics collector
func GetGlobalMetrics() *MetricsCollector {
	return globalMetrics
}