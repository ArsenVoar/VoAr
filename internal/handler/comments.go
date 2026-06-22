package handler

import (
	"VoAr/internal/models"
	"VoAr/internal/service"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (h *Handler) CreateComment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)

	content := r.FormValue("content")

	articleID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "invalid article id", http.StatusBadRequest)
		return
	}

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

	comment := models.Comment{
		ArticleID: articleID,
		UserID:    userID,
		Content:   content,
	}

	_, err = h.CommentService.CreateComment(ctx, comment)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrEmptyFields):
			http.Error(w, "empty fields", http.StatusBadRequest)

		case errors.Is(err, service.ErrNotFound):
			http.Error(w, "post not found", http.StatusNotFound)

		case errors.Is(err, service.ErrInvalidInput):
			http.Error(w, "invalid input", http.StatusBadRequest)

		default:
			log.Printf("internal error: %v", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	http.Redirect(w, r, "/show/"+vars["id"], http.StatusSeeOther)
}
