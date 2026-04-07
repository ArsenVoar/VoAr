package handler

import (
	"VoAr/internal/contextkeys"
	"VoAr/internal/middleware"
	"VoAr/internal/models"
	"VoAr/internal/repository"
	"VoAr/internal/service"
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/markbates/goth/gothic"
)

func HandleFunc(db *sql.DB) *http.Server {
	router := mux.NewRouter()

	router.Use(middleware.DbMiddleware(db))

	router.HandleFunc("/", MainPage).Methods("GET")
	router.HandleFunc("/create", Create).Methods("GET")
	router.HandleFunc("/examples", Examples).Methods("GET")
	router.HandleFunc("/chat", Chat).Methods("GET")
	router.HandleFunc("/userSavedSuccesfull", UserSavedSuccesfull).Methods("GET")
	router.HandleFunc("/userExists", UserExists).Methods("GET")

	router.HandleFunc("/post", Post).Methods("GET")
	router.HandleFunc("/show/{id:[0-9]+}", ShowPost).Methods("GET")
	router.HandleFunc("/profile/{id:[0-9]+}", UserProfile).Methods("GET")

	router.HandleFunc("/save_article", SaveArticle).Methods("POST")

	router.HandleFunc("/auth/google/callback", func(w http.ResponseWriter, r *http.Request) {
		_, err := gothic.CompleteUserAuth(w, r)
		if err != nil {
			log.Printf("Error completing user authentication: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	router.HandleFunc("/auth/google", func(w http.ResponseWriter, r *http.Request) {
		gothic.BeginAuthHandler(w, r)
	})

	router.HandleFunc("/googleSignIn", GoogleSignIn).Methods("GET")

	router.HandleFunc("/save_user", func(w http.ResponseWriter, r *http.Request) {
		ctxDB := r.Context().Value(contextkeys.DbKey).(*sql.DB)

		name := r.FormValue("name")
		email := r.FormValue("email")

		repo := &repository.UserRepository{DB: ctxDB}
		svc := &service.UserService{Repo: repo}

		user := models.User{
			Name:  name,
			Email: email,
		}

		err := svc.Create(user)
		if err != nil {
			if errors.Is(err, service.ErrUserExists) {
				http.Redirect(w, r, "/userExists", http.StatusSeeOther)
				return
			}
			log.Printf("Error saving user: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}).Methods("POST")

	staticFileDirectory := http.Dir("web/css")
	router.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(staticFileDirectory)))

	return &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
}
