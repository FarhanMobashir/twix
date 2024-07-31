package middlewares

import (
	"net/http"
	"time"

	"github.com/fatih/color"
)

// Middleware function for logging with colored output
func LoggingMiddleware(next http.Handler) http.Handler {
	// Define colors
	infoColor := color.New(color.FgGreen).PrintfFunc()
	errorColor := color.New(color.FgRed).PrintfFunc()
	resetColor := color.New(color.Reset).PrintfFunc()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Log the request start
		infoColor("Started %s %s\n", r.Method, r.URL.Path)

		// Create a response writer to capture the status code
		lrw := &loggingResponseWriter{w, http.StatusOK}
		next.ServeHTTP(lrw, r)

		// Log the request completion
		duration := time.Since(start)
		statusColor := infoColor
		if lrw.status >= 500 {
			statusColor = errorColor
		}
		statusColor("Completed %s %s with status %d in %v\n", r.Method, r.URL.Path, lrw.status, duration)

		// Reset color at the end
		resetColor("")
	})
}

// loggingResponseWriter wraps around http.ResponseWriter to capture the status code
type loggingResponseWriter struct {
	http.ResponseWriter
	status int
}

// WriteHeader overrides the default WriteHeader to capture the status code
func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.status = code
	lrw.ResponseWriter.WriteHeader(code)
}

// Write overrides the default Write to support capturing the status code
func (lrw *loggingResponseWriter) Write(b []byte) (int, error) {
	return lrw.ResponseWriter.Write(b)
}
