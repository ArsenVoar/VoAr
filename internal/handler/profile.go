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
	"github.com/markbates/goth/gothic"
)

func UserProfile(w http.ResponseWriter, r *http.Request) {
	ctxDB := r.Context().Value(contextkeys.DbKey).(*sql.DB)

	vars := mux.Vars(r)

	repo := &repository.UserRepository{DB: ctxDB}
	svc := &service.UserService{Repo: repo}

	session, err := gothic.Store.Get(r, "session-name")
	if err != nil {
		log.Printf("Error getting session: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	sessionUserID, ok := session.Values["userId"].(string)
	if !ok {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	user, err := svc.GetProfile(vars["id"], sessionUserID)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		if errors.Is(err, service.ErrUnauthorized) {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		log.Printf("Error getting profile: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	renderTemplate(w, "profile", user)
}
