package handler

import (
	"net/http"
	"time"

	"VoAr/internal/logger"
	"VoAr/internal/middleware"

	"github.com/gorilla/mux"
	"github.com/markbates/goth/gothic"
)

// SetupRouter configures the application routes and middleware.
func SetupRouter(h *Handler) *http.Server {
	router := mux.NewRouter()

	// Global middleware.
	router.Use(middleware.RequestIDMiddleware)
	router.Use(middleware.TimeOutMiddleware)
	router.Use(middleware.LoggingMiddleware)

	authMiddleware := middleware.AuthMiddleware(h.Store)

	api := &ApiHandler{
		UserService: h.UserService,
		PostService: h.PostService,
	}

	// Debug routes.
	router.HandleFunc("/debug/slow", h.SlowHandler)

	// Public pages.
	router.HandleFunc("/", h.MainPage).Methods("GET")
	router.HandleFunc("/auth", h.AuthPage).Methods("GET")
	router.HandleFunc("/create", h.Create).Methods("GET")
	router.HandleFunc("/examples", h.Examples).Methods("GET")
	router.HandleFunc("/post", h.GetPosts).Methods("GET")
	router.HandleFunc("/show/{id:[0-9]+}", h.GetPostByID).Methods("GET")

	// Authentication.
	router.HandleFunc("/register", h.Register).Methods("POST")
	router.HandleFunc("/logout", h.Logout).Methods("GET")
	router.HandleFunc("/login", h.Login).Methods("POST")

	// REST API.
	router.HandleFunc("/api/login", api.Login).Methods("POST")
	router.HandleFunc("/api/posts", api.GetPosts).Methods("GET")

	// Protected pages.
	router.Handle(
		"/notifications",
		authMiddleware(http.HandlerFunc(h.Notifications)),
	).Methods("GET")

	router.Handle(
		"/profile/{id:[0-9]+}",
		authMiddleware(http.HandlerFunc(h.ShowProfile)),
	).Methods("GET")

	router.Handle(
		"/post/{id:[0-9]+}/comments",
		authMiddleware(http.HandlerFunc(h.CreateComment)),
	).Methods("POST")

	router.Handle(
		"/save_post",
		authMiddleware(http.HandlerFunc(h.CreatePost)),
	).Methods("POST")

	// Google OAuth.
	router.HandleFunc("/auth/google/callback", func(w http.ResponseWriter, r *http.Request) {
		_, err := gothic.CompleteUserAuth(w, r)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})
	router.HandleFunc("/auth/google", func(w http.ResponseWriter, r *http.Request) {
		gothic.BeginAuthHandler(w, r)
	})

	// Static files.
	staticFileDirectory := http.Dir("web/css")
	router.PathPrefix("/css/").
		Handler(http.StripPrefix("/css/", http.FileServer(staticFileDirectory)))

	return &http.Server{
		Addr:              ":8080",
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
}
