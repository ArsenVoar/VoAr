package middleware

import (
	"VoAr/internal/contextkeys"
	"context"
	"net/http"
	"strconv"
	"time"
)

func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		requestID := generateID()

		ctx := context.WithValue(
			r.Context(),
			contextkeys.RequestIDKey,
			requestID,
		)

		newRequest := r.WithContext(ctx)

		next.ServeHTTP(w, newRequest)
	})
}

func generateID() string {
	return strconv.FormatInt(time.Now().UnixNano(), 10)
}
