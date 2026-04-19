package logging

import (
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

// RequestIDHeader is the header key for the request ID.
const RequestIDHeader = "X-Request-ID"

// responseWriter is a wrapper around http.ResponseWriter that captures the status code.
type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}
	rw.status = code
	rw.wroteHeader = true
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.wroteHeader {
		rw.WriteHeader(http.StatusOK)
	}
	return rw.ResponseWriter.Write(b)
}

func (rw *responseWriter) Unwrap() http.ResponseWriter {
	return rw.ResponseWriter
}

// LoggingMiddleware logs each incoming HTTP request and injects a logger
// with request metadata and a unique Request ID into the request context.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Get or generate Request ID
		requestID := r.Header.Get(RequestIDHeader)
		if requestID == "" {
			requestID = uuid.Must(uuid.NewV7()).String()
		}

		// Set Request ID in response header
		w.Header().Set(RequestIDHeader, requestID)

		// Create a logger for this specific request with base metadata
		logger := slog.Default().With(
			slog.String("request_id", requestID),
			slog.String("method", r.Method),
			slog.String("remote_addr", r.RemoteAddr),
		)

		// Inject the logger into the context
		ctx := WithLogger(r.Context(), logger)
		r = r.WithContext(ctx)

		// Wrap response writer to capture status
		rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}

		// Call the next handler
		next.ServeHTTP(rw, r)

		// Log completion
		route := r.Pattern
		if i := strings.Index(route, " "); i != -1 {
			route = route[i+1:]
		}

		logger.Info("request completed",
			slog.String("request_route", route),
			slog.Int("status", rw.status),
			slog.Int64("duration_ms", time.Since(start).Milliseconds()),
		)
	})
}
