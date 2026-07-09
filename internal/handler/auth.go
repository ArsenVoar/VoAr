package handler

import (
	"VoAr/internal/logger"
	"VoAr/internal/models"
	"VoAr/internal/service"
	"errors"
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
			logger.Error(err.Error())
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
			logger.Error(err.Error())
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	session, err := h.Store.Get(r, "session-name")
	if err != nil {
		logger.Error(err.Error())

		session, err = h.Store.New(r, "session-name")
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		session.Options.MaxAge = -1

		if err := session.Save(r, w); err != nil {
			logger.Error(err.Error())
		}

		http.Redirect(w, r, "/auth", http.StatusSeeOther)
		return
	}

	session.Values["userId"] = user.ID

	err = session.Save(r, w)
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	session, err := h.Store.Get(r, "session-name")
	if err != nil {
		logger.Error(err.Error())

		session, err = h.Store.New(r, "session-name")
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		session.Options.MaxAge = -1

		if err := session.Save(r, w); err != nil {
			logger.Error(err.Error())
		}

		http.Redirect(w, r, "/auth", http.StatusSeeOther)
		return
	}

	session.Options.MaxAge = -1

	err = session.Save(r, w)
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
