package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

// Logger handles application and audit logging
type Logger struct {
	appLogger   *log.Logger
	auditLogger *log.Logger
	logDir      string
}

// NewLogger creates a new logger instance
func NewLogger(logDir string) (*Logger, error) {
	// Create log directory if not exists
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	// Application log file
	appLogFile, err := openLogFile(filepath.Join(logDir, "app.log"))
	if err != nil {
		return nil, fmt.Errorf("failed to open app log file: %w", err)
	}

	// Audit log file
	auditLogFile, err := openLogFile(filepath.Join(logDir, "audit.log"))
	if err != nil {
		return nil, fmt.Errorf("failed to open audit log file: %w", err)
	}

	return &Logger{
		appLogger:   log.New(appLogFile, "", log.LstdFlags),
		auditLogger: log.New(auditLogFile, "", log.LstdFlags|log.LUTC),
		logDir:      logDir,
	}, nil
}

// openLogFile opens or creates a log file with append mode
func openLogFile(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
}

// Info logs an info message
func (l *Logger) Info(message string) {
	l.appLogger.Printf("[INFO] %s", message)
}

// Error logs an error message
func (l *Logger) Error(message string) {
	l.appLogger.Printf("[ERROR] %s", message)
}

// Warn logs a warning message
func (l *Logger) Warn(message string) {
	l.appLogger.Printf("[WARN] %s", message)
}

// Debug logs a debug message
func (l *Logger) Debug(message string) {
	l.appLogger.Printf("[DEBUG] %s", message)
}

// Audit logs an audit entry
func (l *Logger) Audit(method, endpoint, payload, response, ipAddress string, statusCode int) {
	timestamp := time.Now().UTC().Format(time.RFC3339)

	logEntry := fmt.Sprintf(
		"\n------------------------------------------------------------------------------------------------\n"+
			"timestamp=%s\n"+
			"method=%s\n"+
			"endpoint=%s\n"+
			"payload=%s\n"+
			"status=%d\n"+
			"response=%s\n"+
			"ip=%s"+
			"\n------------------------------------------------------------------------------------------------\n",
		timestamp,
		method,
		endpoint,
		payload,
		statusCode,
		response,
		ipAddress,
	)

	l.auditLogger.Println(logEntry)
}

// AuditJSON logs audit entry in JSON format
func (l *Logger) AuditJSON(data map[string]interface{}) {
	timestamp := time.Now().UTC().Format(time.RFC3339)
	data["timestamp"] = timestamp
	data["type"] = "audit"

	// Simple JSON-like format
	logEntry := "[AUDIT] "
	for key, value := range data {
		logEntry += fmt.Sprintf(`%s=%v `, key, value)
	}

	l.auditLogger.Println(logEntry)
}

// HTTP logs an HTTP request
func (l *Logger) HTTP(method, path, ipAddress string, statusCode int, duration time.Duration, userID *int64) {
	userIDStr := "-"
	if userID != nil {
		userIDStr = fmt.Sprintf("%d", *userID)
	}

	// Apache-style combined log format
	logEntry := fmt.Sprintf(
		`%s - %s [%s] "%s %s" %d %v`,
		ipAddress,
		userIDStr,
		time.Now().Format("02/Jan/2006:15:04:05 -0700"),
		method,
		path,
		statusCode,
		duration,
	)

	l.appLogger.Println(logEntry)
}
