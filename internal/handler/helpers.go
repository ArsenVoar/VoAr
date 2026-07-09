package handler

import (
	"VoAr/internal/logger"
	"encoding/json"
	"errors"
	"net/http"
)

// ErrInvalidSession indicates that the current session is missing or corrupted.
var ErrInvalidSession = errors.New("invalid session")

// currentUserID extracts the authenticated user's ID from the session.
func (h *Handler) currentUserID(w http.ResponseWriter, r *http.Request) (int, error) {
	session, err := h.Store.Get(r, "session-name")
	if err != nil {
		session, err = h.Store.New(r, "session-name")
		if err != nil {
			return 0, ErrInvalidSession
		}

		session.Options.MaxAge = -1

		if err := session.Save(r, w); err != nil {
			return 0, ErrInvalidSession
		}

		return 0, ErrInvalidSession
	}

	userIDValue, ok := session.Values["userId"]
	if !ok {
		return 0, ErrInvalidSession
	}

	userID, ok := userIDValue.(int)
	if !ok {
		return 0, ErrInvalidSession
	}

	return userID, nil
}

// writeJSON writes a JSON response with the provided HTTP status code.
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		logger.Error(err.Error())
	}
}
