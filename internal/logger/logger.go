package logger

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

type LogEntry struct {
	Timestamp    time.Time `json:"timestamp"`
	Method       string    `json:"method"`
	Path         string    `json:"path"`
	StatusCode   int       `json:"status_code"`
	ResponseTime int64     `json:"response_time_ms"`
	UserAgent    string    `json:"user_agent,omitempty"`
	RemoteAddr   string    `json:"remote_addr,omitempty"`
	RequestID    string    `json:"request_id,omitempty"`
	Error        string    `json:"error,omitempty"`
	RequestBody  string    `json:"request_body,omitempty"`
	ResponseSize int       `json:"response_size,omitempty"`
}

type Logger struct {
	file *os.File
}

var APILogger *Logger

// Initialize creates a new logger instance and sets up the log file
func Initialize() error {
	// Create logs directory if it doesn't exist
	logsDir := "logs"
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		return fmt.Errorf("failed to create logs directory: %v", err)
	}

	// Create log file with current date
	logFileName := fmt.Sprintf("api_%s.log", time.Now().Format("2006-01-02"))
	logFilePath := filepath.Join(logsDir, logFileName)

	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("failed to open log file: %v", err)
	}

	APILogger = &Logger{file: file}
	log.Printf("API logging initialized: %s", logFilePath)
	return nil
}

// Close closes the log file
func Close() error {
	if APILogger != nil && APILogger.file != nil {
		return APILogger.file.Close()
	}
	return nil
}

// LogRequest logs an API request with all relevant details
func (l *Logger) LogRequest(entry LogEntry) {
	if l == nil || l.file == nil {
		return
	}

	// Convert to JSON
	jsonData, err := json.Marshal(entry)
	if err != nil {
		log.Printf("Failed to marshal log entry: %v", err)
		return
	}

	// Write to file
	if _, err := l.file.WriteString(string(jsonData) + "\n"); err != nil {
		log.Printf("Failed to write to log file: %v", err)
		return
	}

	// Also log to console for development
	statusColor := getStatusColor(entry.StatusCode)
	fmt.Printf("%s [%s] %s %s - %d (%dms)%s\n",
		entry.Timestamp.Format("2006-01-02 15:04:05"),
		entry.Method,
		entry.Path,
		statusColor,
		entry.StatusCode,
		entry.ResponseTime,
		"\033[0m", // Reset color
	)
}

// getStatusColor returns ANSI color codes based on HTTP status code
func getStatusColor(statusCode int) string {
	switch {
	case statusCode >= 200 && statusCode < 300:
		return "\033[32m" // Green
	case statusCode >= 300 && statusCode < 400:
		return "\033[33m" // Yellow
	case statusCode >= 400 && statusCode < 500:
		return "\033[31m" // Red
	case statusCode >= 500:
		return "\033[35m" // Magenta
	default:
		return "\033[37m" // White
	}
}

// LogError logs an error with context
func (l *Logger) LogError(method, path, errorMsg string) {
	entry := LogEntry{
		Timestamp: time.Now(),
		Method:    method,
		Path:      path,
		Error:     errorMsg,
	}
	l.LogRequest(entry)
}