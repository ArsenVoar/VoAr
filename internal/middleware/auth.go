package middleware

import (
	"net/http"

	"github.com/gorilla/sessions"
)

func AuthMiddleware(store *sessions.CookieStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			session, err := store.Get(r, "session-name")
			if err != nil {
				http.Error(w, "internal server error", http.StatusInternalServerError)
				return
			}

			userID, ok := session.Values["userId"].(int)
			if !ok || userID <= 0 {
				http.Redirect(w, r, "/auth", http.StatusSeeOther)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
