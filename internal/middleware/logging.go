package middleware

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"backend/internal/logger"
)

// ResponseWriter wrapper to capture response data
type responseWriter struct {
	http.ResponseWriter
	statusCode int
	body       *bytes.Buffer
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if rw.body != nil {
		rw.body.Write(b)
	}
	return rw.ResponseWriter.Write(b)
}

// LoggingMiddleware logs all HTTP requests and responses
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Read request body for logging (if it's not too large)
		var requestBody string
		if r.Body != nil && r.ContentLength > 0 && r.ContentLength < 10240 { // Max 10KB
			bodyBytes, err := io.ReadAll(r.Body)
			if err == nil {
				requestBody = string(bodyBytes)
				// Restore the body for the actual handler
				r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			}
		}

		// Create response writer wrapper
		rw := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK, // Default status
			body:           &bytes.Buffer{},
		}

		// Call the next handler
		next.ServeHTTP(rw, r)

		// Calculate response time
		duration := time.Since(start)

		// Create log entry
		entry := logger.LogEntry{
			Timestamp:    start,
			Method:       r.Method,
			Path:         r.URL.Path,
			StatusCode:   rw.statusCode,
			ResponseTime: duration.Milliseconds(),
			UserAgent:    r.UserAgent(),
			RemoteAddr:   r.RemoteAddr,
			RequestBody:  requestBody,
			ResponseSize: rw.body.Len(),
		}

		// Log the request
		if logger.APILogger != nil {
			logger.APILogger.LogRequest(entry)
		}
	})
}