package handler

import (
	"VoAr/internal/logger"
	"VoAr/internal/service"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (h *Handler) ShowProfile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	sessionUserID, err := h.currentUserID(w, r)
	if err != nil {
		return
	}

	ctx := r.Context()

	user, err := h.UserService.GetUserProfile(ctx, userID, sessionUserID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUserNotFound):
			http.Error(w, "user not found", http.StatusNotFound)

		case errors.Is(err, service.ErrUnauthorized):
			http.Redirect(w, r, "/", http.StatusSeeOther)

		default:
			logger.Error(err.Error())
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	h.renderTemplate(w, "profile", user, "", r)
}
