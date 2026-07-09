package handler

import (
	"VoAr/internal/logger"
	"html/template"
	"net/http"
)

var tmpl = template.Must(template.ParseGlob("web/templates/*.html"))

// renderTemplate renders an HTML template with common page data.
func (h *Handler) renderTemplate(w http.ResponseWriter, name string, data any, errMsg string, r *http.Request) {
	// Ignore invalid sessions on public pages and render the user as logged out.
	session, _ := h.Store.Get(r, "session-name")

	var userID int
	var isAuth bool

	if value, ok := session.Values["userId"].(int); ok {
		isAuth = true
		userID = value
	}

	tmplData := TemplateData{
		UserID: userID,
		IsAuth: isAuth,
		Data:   data,
		Error:  errMsg,
	}
	if err := tmpl.ExecuteTemplate(w, name, tmplData); err != nil {
		logger.Error(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) MainPage(w http.ResponseWriter, r *http.Request) {
	h.renderTemplate(w, "mainPage", nil, "", r)
}

func (h *Handler) Examples(w http.ResponseWriter, r *http.Request) {
	h.renderTemplate(w, "examples", nil, "", r)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	h.renderTemplate(w, "create", nil, "", r)
}

func (h *Handler) AuthPage(w http.ResponseWriter, r *http.Request) {
	h.renderTemplate(w, "auth", nil, "", r)
}
