package middleware

import (
	"VoAr/internal/contextkeys"
	"VoAr/internal/logger"
	"net/http"
	"time"
)

// responseWriter captures the HTTP status code for logging.
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode

	rw.ResponseWriter.WriteHeader(statusCode)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rw := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(rw, r)

		duration := time.Since(start)

		requestID := contextkeys.GetRequestID(r.Context())

		logger.Request(requestID, r.Method, r.URL.Path, rw.statusCode, duration)
	})
}
