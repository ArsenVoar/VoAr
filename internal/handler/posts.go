package handler

import (
	"VoAr/internal/contextkeys"
	"VoAr/internal/repository"
	"VoAr/internal/service"
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func SaveArticle(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	anons := r.FormValue("anons")
	fullText := r.FormValue("full_text")

	ctxDB := r.Context().Value(contextkeys.DbKey).(*sql.DB)

	repo := &repository.ArticleRepository{DB: ctxDB}
	svc := &service.ArticleService{Repo: repo}

	err := svc.Create(title, anons, fullText)
	if err != nil {
		if errors.Is(err, service.ErrEmptyFields) {
			http.Error(w, "Please provide all required fields", http.StatusBadRequest)
			return
		}

		if errors.Is(err, service.ErrNoRows) {
			http.Error(w, "Failed to save article", http.StatusInternalServerError)
			return
		}

		log.Printf("Error creating article: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func Post(w http.ResponseWriter, r *http.Request) {
	ctxDB := r.Context().Value(contextkeys.DbKey).(*sql.DB)

	repo := &repository.PostRepository{DB: ctxDB}
	service := &service.PostService{Repo: repo}

	pageParam := r.FormValue("page")

	posts, err := service.GetPosts(pageParam)
	if err != nil {
		log.Printf("Error getting posts: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	renderTemplate(w, "post", posts)
}

func ShowPost(w http.ResponseWriter, r *http.Request) {
	ctxDB := r.Context().Value(contextkeys.DbKey).(*sql.DB)

	vars := mux.Vars(r)

	repo := &repository.ShowPostRepository{DB: ctxDB}
	service := &service.ShowPostService{Repo: repo}

	post, err := service.GetPost(vars["id"])
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			http.Error(w, "Post not found", http.StatusNotFound)
			return
		}

		log.Printf("Error getting post: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	renderTemplate(w, "show", post)
}
