package handler

import (
	"VoAr/internal/models"
	"VoAr/internal/service"

	"github.com/gorilla/sessions"
)

type Handler struct {
	UserService    *service.UserService
	PostService    *service.PostService
	ArticleService *service.ArticleService
	CommentService *service.CommentService
	Store          *sessions.CookieStore
}

func NewHandler(
	userSvc *service.UserService,
	postSvc *service.PostService,
	articleSvc *service.ArticleService,
	commentsSvc *service.CommentService,
	store *sessions.CookieStore,
) *Handler {
	return &Handler{
		UserService:    userSvc,
		PostService:    postSvc,
		ArticleService: articleSvc,
		CommentService: commentsSvc,
		Store:          store,
	}
}

type TemplateData struct {
	UserID int
	IsAuth bool
	Data   interface{}
	Error  string
}
type ShowPostData struct {
	Post     models.Post
	Comments []models.Comment
}
