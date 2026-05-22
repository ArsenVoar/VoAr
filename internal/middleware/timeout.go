package middleware

import (
	"context"
	"net/http"
	"time"
)

func TimeOutMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(
			r.Context(),
			5*time.Second,
		)

		defer cancel()

		newRequest := r.WithContext(ctx)

		next.ServeHTTP(w, newRequest)
	})
}
