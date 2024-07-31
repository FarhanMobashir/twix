package middlewares

import (
	"net/http"
	"time"

	"github.com/fatih/color"
)

// Middleware function for logging with colored output
func LoggingMiddleware(next http.Handler) http.Handler {
	// Define colors with background
	infoColor := color.New(color.FgBlack).PrintfFunc()
	errorColor := color.New(color.FgRed).PrintfFunc()
	resetColor := color.New(color.Reset).PrintfFunc()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lrw := &loggingResponseWriter{w, http.StatusOK}
		next.ServeHTTP(lrw, r)

		duration := time.Since(start)
		statusColor := infoColor
		if lrw.status >= 500 {
			statusColor = errorColor
		}
		statusColor("%s %s with status %d in %v\n", r.Method, r.URL.Path, lrw.status, duration)

		resetColor("")
	})
}

type loggingResponseWriter struct {
	http.ResponseWriter
	status int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.status = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *loggingResponseWriter) Write(b []byte) (int, error) {
	return lrw.ResponseWriter.Write(b)
}
