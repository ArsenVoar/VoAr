package handler

import (
	"VoAr/internal/service"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func (h *Handler) UserProfile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	session, err := h.Store.Get(r, "session-name")
	if err != nil {
		log.Printf("session get error: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	sessionUserID := fmt.Sprintf("%v", session.Values["userId"])

	ctx := r.Context()

	user, err := h.UserService.GetProfile(ctx, vars["id"], sessionUserID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUserNotFound):
			http.Error(w, "user not found", http.StatusNotFound)

		case errors.Is(err, service.ErrUnauthorized):
			http.Redirect(w, r, "/", http.StatusSeeOther)

		default:
			log.Printf("internal error: %v", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	h.renderTemplate(w, "profile", user, "", r)
}
