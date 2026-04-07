package middleware

import (
	"VoAr/internal/contextkeys"
	"context"
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"
)

func DbMiddleware(db *sql.DB) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), contextkeys.DbKey, db)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
