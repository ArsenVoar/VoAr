package handler

import (
	"VoAr/internal/logger"
	"VoAr/internal/models"
	"VoAr/internal/service"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (h *Handler) CreateComment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)

	content := r.FormValue("content")

	postID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "invalid post id", http.StatusBadRequest)
		return
	}

	userID, err := h.currentUserID(w, r)
	if err != nil {
		http.Redirect(w, r, "/auth", http.StatusSeeOther)
		return
	}

	comment := models.Comment{
		PostID:  postID,
		UserID:  userID,
		Content: content,
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
			logger.Error(err.Error())
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	http.Redirect(w, r, "/show/"+vars["id"], http.StatusSeeOther)
}
