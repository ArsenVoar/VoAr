package handler

import (
	"html/template"
	"log"
	"net/http"
)

type ContextKey string

const DbKey ContextKey = "db"

// вспомогательная функция для рендера шаблонов
func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	t, err := template.ParseFiles(
		"web/templates/"+tmpl+".html",
		"web/templates/header.html",
		"web/templates/footer.html",
	)
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = t.ExecuteTemplate(w, tmpl, data)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func MainPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "mainPage", nil)
}

func Examples(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "examples", nil)
}

func Create(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "create", nil)
}

func Chat(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "chat", nil)
}

func UserSavedSuccesfull(w http.ResponseWriter, r *http.Request) {
	log.Println("User saved successfully")
	renderTemplate(w, "userSavedSuccesfull", nil)
}

func UserExists(w http.ResponseWriter, r *http.Request) {
	log.Println("User already exists")
	renderTemplate(w, "userExists", nil)
}

func GoogleSignIn(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "googleSignIn", nil)
}
