package handler

import (
	"VoAr/internal/service"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (h *Handler) SaveArticle(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	anons := r.FormValue("anons")
	fullText := r.FormValue("full_text")

	ctx := r.Context()

	err := h.ArticleService.Create(ctx, title, anons, fullText)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrEmptyFields):
			http.Error(w, "invalid input", http.StatusBadRequest)

		case errors.Is(err, service.ErrNoRows):
			http.Error(w, "failed to save article", http.StatusInternalServerError)

		default:
			log.Printf("internal error: %v", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) Post(w http.ResponseWriter, r *http.Request) {
	pageParam := r.FormValue("page")

	ctx := r.Context()

	posts, err := h.PostService.GetPosts(ctx, pageParam)
	if err != nil {
		log.Printf("internal error: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	h.renderTemplate(w, "post", posts, "", r)
}

func (h *Handler) ShowPost(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("internal error: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	post, err := h.PostService.GetPost(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrNotFound):
			http.Error(w, "post not found", http.StatusNotFound)

		default:
			log.Printf("internal error: %v", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}
	comments, err := h.CommentService.GetCommentsByArticle(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrNotFound):
			http.Error(w, "comments not found", http.StatusNotFound)

		default:
			log.Printf("internal error: %v", err)
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
