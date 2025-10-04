package domain

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

// LogLevel represents the log level
type LogLevel string

const (
	DEBUG LogLevel = "DEBUG"
	INFO  LogLevel = "INFO"
	WARN  LogLevel = "WARN"
	ERROR LogLevel = "ERROR"
)

// LogEntry represents a structured log entry
type LogEntry struct {
	Timestamp   time.Time              `json:"timestamp"`
	Level       LogLevel               `json:"level"`
	Message     string                 `json:"message"`
	Service     string                 `json:"service"`
	RequestID   string                 `json:"request_id,omitempty"`
	UserID      int                    `json:"user_id,omitempty"`
	GroupID     int                    `json:"group_id,omitempty"`
	ParticipantID int                   `json:"participant_id,omitempty"`
	Duration    int64                  `json:"duration_ms,omitempty"`
	Fields      map[string]interface{} `json:"fields,omitempty"`
}

// Logger provides structured logging
type Logger struct {
	service string
	level   LogLevel
}

// NewLogger creates a new logger
func NewLogger(service string, level LogLevel) *Logger {
	return &Logger{
		service: service,
		level:   level,
	}
}

// Debug logs a debug message
func (l *Logger) Debug(message string, fields ...map[string]interface{}) {
	l.log(DEBUG, message, fields...)
}

// Info logs an info message
func (l *Logger) Info(message string, fields ...map[string]interface{}) {
	l.log(INFO, message, fields...)
}

// Warn logs a warning message
func (l *Logger) Warn(message string, fields ...map[string]interface{}) {
	l.log(WARN, message, fields...)
}

// Error logs an error message
func (l *Logger) Error(message string, fields ...map[string]interface{}) {
	l.log(ERROR, message, fields...)
}

// log writes a structured log entry
func (l *Logger) log(level LogLevel, message string, fields ...map[string]interface{}) {
	// Check if we should log this level
	if !l.shouldLog(level) {
		return
	}
	
	entry := LogEntry{
		Timestamp: time.Now(),
		Level:     level,
		Message:   message,
		Service:   l.service,
	}
	
	// Add fields if provided
	if len(fields) > 0 {
		entry.Fields = fields[0]
	}
	
	// Output as JSON
	jsonData, err := json.Marshal(entry)
	if err != nil {
		// Fallback to simple logging
		log.Printf("[%s] %s: %s", level, l.service, message)
		return
	}
	
	fmt.Fprintln(os.Stdout, string(jsonData))
}

// shouldLog determines if a log level should be output
func (l *Logger) shouldLog(level LogLevel) bool {
	levels := map[LogLevel]int{
		DEBUG: 0,
		INFO:  1,
		WARN:  2,
		ERROR: 3,
	}
	
	return levels[level] >= levels[l.level]
}

// WithRequestID adds a request ID to the logger
func (l *Logger) WithRequestID(requestID string) *Logger {
	return &Logger{
		service: l.service,
		level:   l.level,
	}
}

// WithUserID adds a user ID to the logger
func (l *Logger) WithUserID(userID int) *Logger {
	return &Logger{
		service: l.service,
		level:   l.level,
	}
}

// WithGroupID adds a group ID to the logger
func (l *Logger) WithGroupID(groupID int) *Logger {
	return &Logger{
		service: l.service,
		level:   l.level,
	}
}

// WithParticipantID adds a participant ID to the logger
func (l *Logger) WithParticipantID(participantID int) *Logger {
	return &Logger{
		service: l.service,
		level:   l.level,
	}
}

// Global logger instances
var (
	AppLogger      = NewLogger("app", INFO)
	AuthLogger     = NewLogger("auth", INFO)
	GroupLogger    = NewLogger("group", INFO)
	EmailLogger    = NewLogger("email", INFO)
	AssignmentLogger = NewLogger("assignment", INFO)
)