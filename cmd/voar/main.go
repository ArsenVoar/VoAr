package main

import (
	"VoAr/internal/cache"
	"VoAr/internal/database"
	"VoAr/internal/handler"
	"VoAr/internal/logger"
	"VoAr/internal/repository"
	"VoAr/internal/service"
	google "VoAr/pkg/google"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	err = logger.Init()
	if err != nil {
		log.Fatal("Error initializing logger:", err)
	}

	db, err := database.InitDB()
	if err != nil {
		log.Fatal("Error initializing database:", err)
	}
	defer db.Close()

	const cacheTTL = 5 * time.Minute
	redisCache := cache.NewRedisCache("localhost:6379", cacheTTL)

	userRepo := &repository.UserRepository{DB: db}
	postRepo := &repository.PostRepository{DB: db}
	commentRepo := &repository.CommentRepository{DB: db}
	notificationRepo := &repository.NotificationRepository{DB: db}

	userSvc := &service.UserService{UserRepo: userRepo, DB: db}
	postSvc := &service.PostService{PostRepo: postRepo, Cache: redisCache}
	notificationSvc := &service.NotificationService{NotificationRepo: notificationRepo}
	commentSvc := &service.CommentService{CommentRepo: commentRepo, PostReader: postRepo, Notifier: notificationSvc}

	sessionSecret := os.Getenv("SESSION_SECRET")
	if sessionSecret == "" {
		log.Fatal("SESSION_SECRET is required")
	}

	store := sessions.NewCookieStore(
		[]byte(sessionSecret),
	)
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600 * 8,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}

	h := handler.NewHandler(userSvc, postSvc, commentSvc, notificationSvc, store)

	google.Google()

	server := handler.SetupRouter(h)

	logger.Info("starting VoAr server")

	quit := make(chan os.Signal, 1)
	signal.Notify(
		quit,
		os.Interrupt,
		syscall.SIGTERM,
	)

	go func() {
		err := server.ListenAndServe()

		if err != nil && err != http.ErrServerClosed {
			logger.Error(err.Error())
		}
	}()

	<-quit
	logger.Info("shutdown signal received")

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		logger.Error(err.Error())
	} else {
		logger.Info("VoAr server stopped")
	}
}
