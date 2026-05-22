package handler

import (
	"html/template"
	"net/http"
	"strconv"
	"time"
)

var tmpl = template.Must(template.ParseGlob("web/templates/*.html"))

func (h *Handler) renderTemplate(w http.ResponseWriter, name string, data interface{}, errMsg string, r *http.Request) {
	session, _ := h.Store.Get(r, "session-name")

	val := session.Values["userId"]

	var userID int
	var isAuth bool

	if val != nil {
		isAuth = true

		switch v := val.(type) {
		case int:
			userID = v
		case float64:
			userID = int(v)
		case string:
			id, err := strconv.Atoi(v)
			if err == nil {
				userID = id
			}
		}
	}

	tmplData := TemplateData{
		UserID: userID,
		IsAuth: isAuth,
		Data:   data,
		Error:  errMsg,
	}
	tmpl.ExecuteTemplate(w, name, tmplData)
}

func (h *Handler) MainPage(w http.ResponseWriter, r *http.Request) {
	h.renderTemplate(w, "mainPage", nil, "", r)

	select {
	case <-time.After(10 * time.Second):
		w.Write([]byte("slow response"))

	case <-r.Context().Done():
		http.Error(w, r.Context().Err().Error(), http.StatusRequestTimeout)
		return
	}
}

func (h *Handler) Examples(w http.ResponseWriter, r *http.Request) {
	h.renderTemplate(w, "examples", nil, "", r)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	h.renderTemplate(w, "create", nil, "", r)
}

func (h *Handler) Chat(w http.ResponseWriter, r *http.Request) {
	h.renderTemplate(w, "chat", nil, "", r)
}

func (h *Handler) AuthPage(w http.ResponseWriter, r *http.Request) {
	h.renderTemplate(w, "auth", nil, "", r)
}
