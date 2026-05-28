package handler

import (
	"VoAr/internal/models"
	"VoAr/internal/service"
	"encoding/json"
	"errors"
	"net/http"
)

type ApiHandler struct {
	UserService *service.UserService
	PostService *service.PostService
}
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type LoginResponse struct {
	UserID string `json:"userId,omitempty"`
	Error  string `json:"error,omitempty"`
}

func (h *ApiHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(LoginResponse{
			Error: "invalid request",
		})
		return
	}

	ctx := r.Context()

	user, err := h.UserService.Login(ctx, req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidInput):
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(LoginResponse{
				Error: "invalid input",
			})

		case errors.Is(err, service.ErrUnauthorized):
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(LoginResponse{
				Error: "invalid credentials",
			})

		default:
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(LoginResponse{
				Error: "internal error",
			})
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(LoginResponse{
		UserID: user.Id,
	})
}

func (h *ApiHandler) GetPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx := r.Context()

	posts, err := h.PostService.GetPosts(ctx, r.URL.Query().Get("page"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data":  nil,
			"error": "internal error",
		})
		return
	}

	if posts == nil {
		posts = []models.Post{}
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":  posts,
		"error": nil,
	})
}
