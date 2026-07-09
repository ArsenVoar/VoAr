package handler

import (
	"VoAr/internal/models"
	"VoAr/internal/service"
	"encoding/json"
	"errors"
	"net/http"
)

// ApiHandler serves JSON API endpoints.
type ApiHandler struct {
	UserService *service.UserService
	PostService *service.PostService
}

// LoginRequest represents API login credentials.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse is returned by the login endpoint.
type LoginResponse struct {
	UserID int    `json:"userId,omitempty"`
	Error  string `json:"error,omitempty"`
}

func (h *ApiHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, LoginResponse{Error: "invalid request"})

		return
	}

	ctx := r.Context()

	user, err := h.UserService.Login(ctx, req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidInput):
			writeJSON(w, http.StatusBadRequest, LoginResponse{Error: "invalid input"})

		case errors.Is(err, service.ErrUnauthorized):
			writeJSON(w, http.StatusUnauthorized, LoginResponse{Error: "invalid credentials"})

		default:
			writeJSON(w, http.StatusInternalServerError, LoginResponse{Error: "internal error"})
		}
		return
	}

	writeJSON(w, http.StatusOK, LoginResponse{
		UserID: user.ID,
	})
}

func (h *ApiHandler) GetPosts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	posts, err := h.PostService.GetPosts(ctx, r.URL.Query().Get("page"))
	if errors.Is(err, service.ErrInvalidInput) {
		writeJSON(w, http.StatusBadRequest, map[string]any{"data": nil, "error": "invalid input"})
		return
	}
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"data": nil, "error": "internal error"})
		return
	}
	if posts == nil {
		posts = []models.Post{}
	}

	writeJSON(w, http.StatusOK, map[string]any{"data": posts, "error": nil})
}
