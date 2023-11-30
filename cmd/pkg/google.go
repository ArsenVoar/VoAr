// Package google provides functions for configuring and handling Google authentication
package google

import (
	"log"           //Package log implements a simple logging package
	"net/http"      //Package http provides HTTP client and server implementations
	"os"            //Package os provides a way of using operating system-dependent functionality
	"text/template" //Package template 	implements data-driven templates for generating textual output

	"github.com/gorilla/pat"                     // Package pat implements a request router and dispatcher
	"github.com/gorilla/sessions"                // Package sessions provides cookie and filesystem sessions and infrastructure for custom session backends
	"github.com/joho/godotenv"                   // Package godotenv loads environment variables from a file
	"github.com/markbates/goth"                  // Package goth provides a simple, clean, and idiomatic way to write authentication packages
	"github.com/markbates/goth/gothic"           // Package gothic provides the ability to use multiple providers for authentication
	"github.com/markbates/goth/providers/google" // Package google implements the OAuth2 protocol for authenticating users with Google
)

// Google function configures and sets up Google authentication using the Goth package.
func Google() {
	// Loading environment variables from the specified file
	err := godotenv.Load("C:/Arsen/VSC/Projects/VoAr/st.env")
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	// Retrieving configuration values from environment variables
	sessionKey := os.Getenv("SESSION_KEY")
	maxAge := 86400 * 30
	isProd := false

	//Configuring the session store for cookie-based sessions
	store := sessions.NewCookieStore([]byte(sessionKey))
	store.MaxAge(maxAge)
	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = isProd

	//Setting the configured session store as the store used by Goth
	gothic.Store = store

	//Retrieving Google 0Auth client credentials from environment variables
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	callbackURL := os.Getenv("GOOGLE_CALLBACK_URL")

	//Configuring and adding Google as an authentication provider to Goth
	goth.UseProviders(
		google.New(clientID, clientSecret, callbackURL, "email", "profile"),
	)

	//Creating a new router from the Gorrila Pat package
	p := pat.New()

	//Handling the callback URL for the Google authentication provider
	p.Get("/auth/{provider/callback", func(res http.ResponseWriter, req *http.Request) {
		//Completing the user authentication process
		user, err := gothic.CompleteUserAuth(res, req)
		if err != nil {
			http.Error(res, "Internal server error", http.StatusInternalServerError)
			log.Printf("Error getting provider name: %v", err)
			return
		}

		//Storing user information in the session
		session, _ := gothic.Store.Get(req, "session-name")
		session.Values["user"] = user
		session.Save(req, res)

		// Redirecting to the complete page after successful authentication
		http.Redirect(res, req, "/complete", http.StatusSeeOther)
	})

	//Handling the initiation of Google authentication
	p.Get("/auth/{provider}", func(res http.ResponseWriter, req *http.Request) {
		gothic.BeginAuthHandler(res, req)
	})

	//Handling the root URL, rendering the main page template
	p.Get("/", func(res http.ResponseWriter, req *http.Request) {
		t, err := template.ParseFiles("templates/mainPage.html")
		if err != nil {
			http.Error(res, "Internal server error", http.StatusInternalServerError)
			log.Printf("Error parsing template: %v", err)
			return
		}
		t.Execute(res, false)
	})
}
