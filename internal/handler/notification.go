package handler

import (
	"log"
	"net/http"
)

func (h *Handler) Notifications(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	session, err := h.Store.Get(r, "session-name")
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	userIDValue, ok := session.Values["userId"]
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	userID, ok := userIDValue.(int)
	if !ok {
		http.Error(w, "invalid user id", http.StatusInternalServerError)
		return
	}

	notifications, err := h.NotificationService.GetNotificationsByUser(ctx, userID)
	if err != nil {
		log.Printf("internal error: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	h.renderTemplate(w, "notifications", notifications, "", r)
}
