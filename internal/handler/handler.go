package handler

import (
	"VoAr/internal/service"

	"github.com/gorilla/sessions"
)

type Handler struct {
	UserService    *service.UserService
	PostService    *service.PostService
	ArticleService *service.ArticleService
	Store          *sessions.CookieStore
}

func NewHandler(
	userSvc *service.UserService,
	postSvc *service.PostService,
	articleSvc *service.ArticleService,
	store *sessions.CookieStore,
) *Handler {
	return &Handler{
		UserService:    userSvc,
		PostService:    postSvc,
		ArticleService: articleSvc,
		Store:          store,
	}
}

type TemplateData struct {
	UserID int
	IsAuth bool
	Data   interface{}
	Error  string
}
