package handler

import (
	"VoAr/internal/logger"
	"VoAr/internal/service"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (h *Handler) CreatePost(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	anons := r.FormValue("anons")
	fullText := r.FormValue("full_text")

	userID, err := h.currentUserID(w, r)
	if err != nil {
		http.Redirect(w, r, "/auth", http.StatusSeeOther)
		return
	}

	ctx := r.Context()

	err = h.PostService.CreatePost(ctx, userID, title, anons, fullText)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrEmptyFields):
			http.Error(w, "invalid input", http.StatusBadRequest)

		case errors.Is(err, service.ErrNoRowsAffected):
			http.Error(w, "failed to save post", http.StatusInternalServerError)

		default:
			logger.Error(err.Error())
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) GetPosts(w http.ResponseWriter, r *http.Request) {
	pageParam := r.FormValue("page")

	ctx := r.Context()

	posts, err := h.PostService.GetPosts(ctx, pageParam)
	if errors.Is(err, service.ErrInvalidInput) {
		http.Error(w, "error invalid input", http.StatusBadRequest)
		return
	}
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	h.renderTemplate(w, "post", posts, "", r)
}

func (h *Handler) GetPostByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "invalid post id", http.StatusBadRequest)
		return
	}
	post, err := h.PostService.GetPostByID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrNotFound):
			http.Error(w, "post not found", http.StatusNotFound)

		default:
			logger.Error(err.Error())
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}
	comments, err := h.CommentService.GetCommentsByPost(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrNotFound):
			http.Error(w, "comments not found", http.StatusNotFound)

		default:
			logger.Error(err.Error())
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	data := ShowPostData{
		Post:     post,
		Comments: comments,
	}

	h.renderTemplate(w, "show", data, "", r)
}
