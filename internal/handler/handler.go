package handler

import (
	"VoAr/internal/models"
	"VoAr/internal/service"

	"github.com/gorilla/sessions"
)

// Handler contains dependencies shared across HTTP handlers.
type Handler struct {
	UserService         *service.UserService
	PostService         *service.PostService
	CommentService      *service.CommentService
	NotificationService *service.NotificationService
	Store               *sessions.CookieStore
}

func NewHandler(
	userSvc *service.UserService,
	postSvc *service.PostService,
	commentsSvc *service.CommentService,
	notificationService *service.NotificationService,
	store *sessions.CookieStore,
) *Handler {
	return &Handler{
		UserService:         userSvc,
		PostService:         postSvc,
		CommentService:      commentsSvc,
		NotificationService: notificationService,
		Store:               store,
	}
}

// TemplateData is passed to HTML templates.
type TemplateData struct {
	UserID int
	IsAuth bool
	Data   interface{}
	Error  string
}

// ShowPostData combines a post with its comments for rendering.
type ShowPostData struct {
	Post     models.Post
	Comments []models.Comment
}
