package app

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"           // Package for HTTP request multiplexer (router)
	"github.com/markbates/goth"        // Package for multi-provider authentication
	"github.com/markbates/goth/gothic" // Package for handling OAuth state and session
)

// initDB initializes the database connection and performs the necessary checks.
// It reads the required environment variables for database configuration, establishes a connection,
// And checks the connection by pinging the database. If successful, it returns a pointer to the
// Established database connection.
func InitDB() (*sql.DB, error) {
	// Retrieve database configuration parameters from environment variables
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	// Construct connection string using retrieved parameters
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=require host=127.0.0.1 port=5432", user, password, dbname)

	// Open a connection to the PostgreSQL database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Check the connection to the database by pinging it
	if err = db.Ping(); err != nil {
		db.Close() // Close the connection in case of an error
		return nil, err
	}

	// Return the established database connection and nil error
	return db, nil
}

// HandleFunc creates a new HTTP server instance with the specified database connection
// And sets up the routing for various endpoints using the Gorilla Mux router
// It includes middleware to inject the database into the request context
// The function returns the configured HTTP server
func HandleFunc(db *sql.DB) *http.Server {
	//Creating 	a new Gorilla Mux router
	router := mux.NewRouter()

	//Using the dbMiddleware to inject the database into the request context
	router.Use(dbMiddleware(db, router))

	//Handling different routes with corresponding HTTP methods
	router.HandleFunc("/", mainPage).Methods("GET")
	router.HandleFunc("/create", create).Methods("GET")
	router.HandleFunc("/complete", complete).Methods("GET")
	router.HandleFunc("/examples", examples).Methods("GET")
	router.HandleFunc("/chat", chat).Methods("GET")

	//Handling the "/post" endpoint with the post function
	router.HandleFunc("/post", func(w http.ResponseWriter, r *http.Request) {
		//Retrieving the database connection from the request context
		db = r.Context().Value(dbKey).(*sql.DB)
		post(w, r, db)
	}).Methods("GET")

	//Handling the "/show/{id:{0-9}+}" endpoint with the showPost function
	router.HandleFunc("/show/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		//Retrieving the database connection from the request context
		db = r.Context().Value(dbKey).(*sql.DB)
		showPost(w, r, db)
	}).Methods("GET")

	//Handling the "/save_article" endpoint with the save_article function
	router.HandleFunc("/save_article", func(w http.ResponseWriter, r *http.Request) {
		//Retrieving the database connection from the request context
		db = r.Context().Value(dbKey).(*sql.DB)
		save_article(w, r)
	}).Methods("POST")

	//Handling 	authentication using third-party providers (0Auth)
	router.HandleFunc("/auth/{provider}", func(w http.ResponseWriter, r *http.Request) {
		gothic.BeginAuthHandler(w, r)
	})

	//Handling the callback from the third-party authentication provider
	router.HandleFunc("/auth/{provider}/callback", func(w http.ResponseWriter, r *http.Request) {
		user, err := gothic.CompleteUserAuth(w, r)
		if err != nil {
			// Handling authentication error by returning an internal server error response
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			log.Printf("Error completing user authentication: %v", err)
			return
		}
		//Redirecting to the completion page with user information
		http.Redirect(w, r, "/complete?name="+user.Name+"&email="+user.Email, http.StatusSeeOther)
	})

	//Handling the "/googleSignIn" endpoit for Google-Sing-In
	router.HandleFunc("/googleSignIn", func(w http.ResponseWriter, r *http.Request) {
		// Parsing HTML template files for the Google Sign-In page, header, and footer
		t, err := template.ParseFiles("cmd/pkg/app/templates/googleSignIn.html", "cmd/pkg/app/templates/header.html", "cmd/pkg/app/templates/footer.html")
		if err != nil {
			// Handling template parsing error by returning an internal server error response
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			log.Printf("Error parsing template files: %v", err)
			return
		}

		//Handling the completion page after Google-Sign-In
		router.HandleFunc("/complete", func(w http.ResponseWriter, r *http.Request) {
			//Retrieving the database connection from the request context
			db = r.Context().Value("db").(*sql.DB)

			//Retrieving user information from the session
			session, _ := gothic.Store.Get(r, "session-name")
			user, ok := session.Values["user"].(goth.User)
			if !ok {
				// Handling the case where the user is not found in the session
				http.Error(w, "User not found", http.StatusNotFound)
				log.Printf("User not found in session")
				return
			}

			// Parsing HTML template files for the completion page, header, and footer
			t, err := template.ParseFiles("cmd/pkg/app/templates/complete.html", "cmd/pkg/app/templates/header.html", "cmd/pkg/app/templates/footer.html")
			if err != nil {
				// Handling template parsing error by returning an internal server error response
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				log.Printf("Error parsing template files: %v", err)
				return
			}

			// Executing the completion template with user information and writing the output to the response writer
			t.ExecuteTemplate(w, "complete", user)
		}).Methods("GET")
		//Executing  the Google-Sign-In template and writing the output to the response writer
		t.ExecuteTemplate(w, "googleSignIn", nil)
	}).Methods("GET")

	// Serving static files from the "/static/" directory
	staticFileDirectory := http.Dir("./static/")
	staticFileHandler := http.StripPrefix("/static/", http.FileServer(staticFileDirectory))
	router.PathPrefix("/static/").Handler(staticFileHandler)

	//Creating an HTTP server instance with the configured router
	server := &http.Server{
		Addr:    ":8080", //Server address and port
		Handler: router,  //Router handling the request
	}

	//Returning the configured HTTP server
	return server
}
