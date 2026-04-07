package google

import (
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/gorilla/pat"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

func Google() {
	err := godotenv.Load("st.env")
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	sessionKey := os.Getenv("SESSION_KEY")
	maxAge := 86400 * 30
	isProd := false

	store := sessions.NewCookieStore([]byte(sessionKey))
	store.MaxAge(maxAge)
	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = isProd

	gothic.Store = store

	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	callbackURL := os.Getenv("GOOGLE_CALLBACK_URL")

	goth.UseProviders(
		google.New(clientID, clientSecret, callbackURL, "email", "profile"),
	)

	p := pat.New()

	p.Get("/auth/{provider}/callback", func(res http.ResponseWriter, req *http.Request) {
		user, err := gothic.CompleteUserAuth(res, req)
		if err != nil {
			http.Error(res, "Internal server error", http.StatusInternalServerError)
			log.Printf("Error getting provider name: %v", err)
			return
		}

		log.Printf("User ID from provider: %v", user.UserID)

		userId := user.UserID
		session, _ := gothic.Store.Get(req, "session-name")
		session.Values["userId"] = userId
		session.Save(req, res)

		log.Printf("User ID saved in session: %v", userId)

		http.Redirect(res, req, "/profile/"+userId, http.StatusSeeOther)
	})

	p.Get("/auth/{provider}", func(res http.ResponseWriter, req *http.Request) {
		gothic.BeginAuthHandler(res, req)
	})

	p.Get("/", func(res http.ResponseWriter, req *http.Request) {
		t, err := template.ParseFiles("web/templates/mainPage.html")
		if err != nil {
			http.Error(res, "Internal server error", http.StatusInternalServerError)
			log.Printf("Error parsing template: %v", err)
			return
		}
		t.Execute(res, false)
	})
}
