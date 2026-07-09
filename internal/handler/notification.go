package handler

import (
	"VoAr/internal/logger"
	"net/http"
)

func (h *Handler) Notifications(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, err := h.currentUserID(w, r)
	if err != nil {
		http.Redirect(w, r, "/auth", http.StatusSeeOther)
		return
	}

	notifications, err := h.NotificationService.GetNotificationsByUser(ctx, userID)
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	h.renderTemplate(w, "notifications", notifications, "", r)
}
