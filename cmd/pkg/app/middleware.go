package app

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"
)

// dbMiddleware is middleware that injects the database instance into the request context
// It takes the database connection, and the next HTTP handler as input, and returns a MiddlewareFunc
func dbMiddleware(db *sql.DB, next http.Handler) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		//Wrapping the next	handler with a custom logic to inject the database into the context
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			//Creating a new context with the database instance and adding it to the request context
			ctx := context.WithValue(r.Context(), dbKey, db)
			//Serving the next HTTP handler with the modified request context
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
