package handler

import (
	"VoAr/internal/models"
	"VoAr/internal/service"
	"errors"
	"log"
	"net/http"
)

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")
	name := r.FormValue("name")

	user := models.User{
		Email:    email,
		Password: password,
		Name:     name,
	}

	ctx := r.Context()

	err := h.UserService.Register(ctx, user)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidInput):
			http.Error(w, "invalid input", http.StatusBadRequest)

		case errors.Is(err, service.ErrUserExists):
			http.Error(w, "user already exists", http.StatusConflict)

		default:
			log.Printf("internal error: %v", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	ctx := r.Context()

	user, err := h.UserService.Login(ctx, email, password)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidInput):
			h.renderTemplate(w, "auth", nil, "Invalid input", r)

		case errors.Is(err, service.ErrUnauthorized):
			h.renderTemplate(w, "auth", nil, "Invalid credentials", r)

		default:
			log.Printf("internal error: %v", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	session, _ := h.Store.Get(r, "session-name")

	session.Values["userId"] = user.ID

	err = session.Save(r, w)
	if err != nil {
		log.Printf("session save error: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	session, err := h.Store.Get(r, "session-name")
	if err != nil {
		log.Printf("session save error: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	session.Options.MaxAge = -1

	session.Save(r, w)
	if err != nil {
		log.Printf("session save error: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
