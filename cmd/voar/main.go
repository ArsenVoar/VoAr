package main

import (
	"VoAr/internal/database"
	"VoAr/internal/handler"
	"VoAr/internal/repository"
	"VoAr/internal/service"
	google "VoAr/pkg/google"
	"log"
	"os"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	db, err := database.InitDB()
	if err != nil {
		log.Fatal("Error initializing database:", err)
	}
	defer db.Close()

	userRepo := &repository.UserRepository{DB: db}
	postRepo := &repository.PostRepository{DB: db}
	articleRepo := &repository.ArticleRepository{DB: db}

	userSvc := &service.UserService{Repo: userRepo}
	postSvc := &service.PostService{Repo: postRepo}
	articleSvc := &service.ArticleService{Repo: articleRepo}

	sessionSecret := os.Getenv("SESSION_SECRET")

	store := sessions.NewCookieStore(
		[]byte(sessionSecret),
	)
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600 * 8,
		HttpOnly: true,
		Secure:   false,
	}

	h := handler.NewHandler(userSvc, postSvc, articleSvc, store)

	google.Google()

	server := handler.SetupRouter(h)

	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Server error:", err)
	}
}
