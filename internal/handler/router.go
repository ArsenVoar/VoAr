package handler

import (
	"log"
	"net/http"

	"VoAr/internal/middleware"

	"github.com/gorilla/mux"
	"github.com/markbates/goth/gothic"
)

func SetupRouter(h *Handler) *http.Server {
	router := mux.NewRouter()

	router.Use(middleware.RequestIDMiddleware)
	router.Use(middleware.TimeOutMiddleware)
	router.Use(middleware.LoggingMiddleware)

	authMiddleware := middleware.AuthMiddleware(h.Store)

	api := &ApiHandler{
		UserService: h.UserService,
		PostService: h.PostService,
	}

	router.HandleFunc("/debug/slow", h.SlowHandler)

	router.HandleFunc("/", h.MainPage).Methods("GET")
	router.HandleFunc("/auth", h.AuthPage).Methods("GET")
	router.HandleFunc("/create", h.Create).Methods("GET")
	router.HandleFunc("/examples", h.Examples).Methods("GET")
	router.HandleFunc("/chat", h.Chat).Methods("GET")
	router.HandleFunc("/logout", h.Logout).Methods("GET")
	router.HandleFunc("/api/posts", api.GetPosts).Methods("GET")
	router.HandleFunc("/register", h.Register).Methods("POST")
	router.HandleFunc("/login", h.Login).Methods("POST")
	router.HandleFunc("/api/login", api.Login).Methods("POST")

	router.HandleFunc("/post", h.Post).Methods("GET")
	router.HandleFunc("/show/{id:[0-9]+}", h.ShowPost).Methods("GET")

	router.Handle(
		"/profile/{id:[0-9]+}",
		authMiddleware(http.HandlerFunc(h.UserProfile)),
	).Methods("GET")
	router.Handle(
		"/article/{id:[0-9]+}/comments",
		authMiddleware(http.HandlerFunc(h.CreateComment)),
	).Methods("POST")

	router.HandleFunc("/save_article", h.SaveArticle).Methods("POST")

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

	staticFileDirectory := http.Dir("web/css")
	router.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(staticFileDirectory)))

	return &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
}
