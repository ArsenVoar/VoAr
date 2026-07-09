package google

import (
	"VoAr/internal/logger"
	"net/http"
	"os"
	"strconv"
	"text/template"

	"github.com/gorilla/pat"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

func Google() {
	sessionKey := os.Getenv("SESSION_SECRET")
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
			logger.Error(err.Error())
			return
		}

		logger.Info("User ID from provider: " + user.UserID)

		// TODO: Map the Google provider ID to the application's users.id.
		userID, err := strconv.Atoi(user.UserID)
		if err != nil {
			http.Error(res, "Invalid user ID from provider", http.StatusInternalServerError)
			logger.Error(err.Error())
			return
		}
		session, _ := gothic.Store.Get(req, "session-name")
		session.Values["userId"] = userID

		if err := session.Save(req, res); err != nil {
			http.Error(res, "Internal server error", http.StatusInternalServerError)
			logger.Error(err.Error())
			return
		}

		logger.Info("User ID saved in session: " + strconv.Itoa(userID))

		http.Redirect(res, req, "/profile/"+strconv.Itoa(userID), http.StatusSeeOther)
	})

	p.Get("/auth/{provider}", func(res http.ResponseWriter, req *http.Request) {
		gothic.BeginAuthHandler(res, req)
	})

	p.Get("/", func(res http.ResponseWriter, req *http.Request) {
		t, err := template.ParseFiles("web/templates/mainPage.html")
		if err != nil {
			http.Error(res, "Internal server error", http.StatusInternalServerError)
			logger.Error(err.Error())
			return
		}
		if err := t.Execute(res, false); err != nil {
			http.Error(res, "Internal server error", http.StatusInternalServerError)
			logger.Error(err.Error())
			return
		}
	})
}
