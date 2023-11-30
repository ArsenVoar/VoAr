package main

import (
	google "VoAr/cmd/pkg" //Importing the Google package for authentication
	"VoAr/cmd/pkg/app"    //importing the app package (assuming it handles application-specific logic)
	"context"             //Package for working with contex in Go
	"database/sql"        //SQL package for working with databases
	"fmt"
	"log"      //Package for logging
	"net/http" //Package for building HTTP servers
	"os"
	"strconv"       //Package for string conversions
	"text/template" //Package for parsing and executing HTML templates

	"github.com/gorilla/mux" //Package for building powerful and flexible HTTP routers
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"              //PostgreSQL driver for the database/sql package
	"github.com/markbates/goth"        //Package for
	"github.com/markbates/goth/gothic" //Package for handling 0Auth authentication with state
)

// Post represents a structure for storing arcticle data
type Post struct {
	Id        int    //Unique identifier for the arcticle
	Title     string //Title of the arcticle
	Anons     string //Brief summary or announcement of the article
	Full_Text string //Full text content of the article
}

// contextKey is a custom type for convience when using it in yhe context
type contextKey string

// dbkey is the context key for passing the database instance to handlers
const dbKey contextKey = "db"

var (
	posts     = []Post{} //Global slice to store multiple article instances
	showItems = Post{}   //Global varriable to store a single article instance for display
)

// initDB initializes the database connection and performs the necessary checks.
// It reads the required environment variables for database configuration, establishes a connection,
// And checks the connection by pinging the database. If successful, it returns a pointer to the
// Established database connection.
func initDB() (*sql.DB, error) {
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

// dbMiddleware is middleware that injects the database instance into the request context
// It takes the database connection, and the next HTTP handler as input, and returns a MiddlewareFunc
func dbMiddleware(db *sql.DB, next http.Handler) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		//Wrapping the next	handler with a custom logic to inject the database into the context
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			//Creating a new context with the database instance and adding it to the request context
			ctx := context.WithValue(r.Context(), dbKey, db)
			//Serving the next HTTP handler with the modified request context
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// mainPage is an HTTP handler function for serving the main page.
// It parses the template files for the main page, header, and footer,
// Handles potential errors, and executes the main page template
func mainPage(w http.ResponseWriter, r *http.Request) {
	// Parsing HTML template files for the main page, header, and footer
	t, err := template.ParseFiles("templates/mainPage.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		// Handling template parsing error by returning an internal server error response
		http.Error(w, "Error interesting data", http.StatusInternalServerError)
		log.Printf("Error scanning data into database: %v", err)
		return
	}

	// Executing the main page template and writing the output to the response writer
	t.ExecuteTemplate(w, "mainPage", nil)
}

// examples is an HTTP handler function for serving the examples page.
// It parses the template files for the examples page, header, and footer,
// Handles potential errors, and executes the examples page template.
func examples(w http.ResponseWriter, r *http.Request) {
	// Parsing HTML template files for the examples page, header, and footer
	t, err := template.ParseFiles("templates/examples.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		// Handling template parsing error by returning an internal server error response
		http.Error(w, "Error interesting data", http.StatusInternalServerError)
		log.Printf("Error scanning data into database: %v", err)
		return
	}
	// Executing the examples page template and writing the output to the response writer
	t.ExecuteTemplate(w, "examples", nil)
}

// create is an HTTP handler function for serving the create page.
// It parses the template files for the create page, header, and footer,
// Handles potential errors, and executes the create page template.
func create(w http.ResponseWriter, r *http.Request) {
	// Parsing HTML template files for the create page, header, and footer
	t, err := template.ParseFiles("templates/create.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		// Handling template parsing error by returning an internal server error response
		http.Error(w, "Error interesting data", http.StatusInternalServerError)
		log.Printf("Error scanning data into database: %v", err)
		return
	}

	// Executing the create page template and writing the output to the response writer
	t.ExecuteTemplate(w, "create", nil)
}

// save_article is an HTTP handler function for saving an article to the database
// It retrieves form values from the request, validates them, and inserts the data into the database
func save_article(w http.ResponseWriter, r *http.Request) {
	//Retrieving form values from the request
	title := r.FormValue("title")
	anons := r.FormValue("anons")
	full_text := r.FormValue("full_text")

	//Validating if all  required fields are provided
	if title == "" || anons == "" || full_text == "" {
		http.Error(w, "Please provide all required fields", http.StatusBadRequest)
		return
	}

	//Retrieving the database instance from the request context
	db := r.Context().Value(dbKey).(*sql.DB)

	//Executing the SQL query to insert the article into the database and getting the result
	result, err := db.Exec("INSERT INTO articles (title, anons, full_text) VALUES ($1, $2, $3) RETURNING id", title, anons, full_text)
	if err != nil {
		//Handling database insertion error by returning an internal server error response
		http.Error(w, "Error interesting data", http.StatusInternalServerError)
		log.Printf("Error scanning data into database: %v", err)
		return
	}

	//Getting the number of rows affected by the insertion
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		//Handling error while getting the number of rows affected
		http.Error(w, "Error getting rows affected", http.StatusInternalServerError)
		log.Printf("Error scanning data into database: %v", err)
		return
	}

	//Checking if any rows 	were affected, if not, returning an internal server error
	if rowsAffected == 0 {
		http.Error(w, "No rows affected data not inserted", http.StatusInternalServerError)
		log.Printf("No rows affected, data not inserted")
		return
	}

	//Redirecting the user to the main page after succesful article inserion
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// post is an HTTP handler function for displaying a list of articles.
// It retrieves a page number from the request parameters, queries the database for articles,
// and renders the list using the post template.
func post(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Parsing HTML template files for the post page, show page, header, and footer
	t, err := template.ParseFiles("templates/post.html", "templates/show.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		// Handling template parsing error by returning an internal server error response
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Printf("Error parsing template files: %v", err)
		return
	}

	//Default page and page size values
	page := 1
	pageSize := 10

	//Checking if a specific page is requested and updating the page variable
	pageParam := r.FormValue("page")
	if pageParam != "" {
		page, err = strconv.Atoi(pageParam)
		if err != nil {
			// Handling error when converting the page parameter to an integer
			http.Error(w, "Invalid page parameter", http.StatusBadRequest)
			log.Printf("Error converting page parameter to integer: %v", err)
			return
		}
	}

	//Calculating the offset for pagination
	offset := (page - 1) * pageSize

	//Querying the database for a list of articles with pagination
	res, err := db.Query("SELECT * FROM articles LIMIT $1 OFFSET $2", pageSize, offset)
	if err != nil {
		// Handling database query error by returning an internal server error response
		http.Error(w, "Internal server eror", http.StatusInternalServerError)
		log.Printf("Error executing database query: %v", err)
		return
	}

	//Creating a slice to store the retrieved articles
	posts = []Post{}
	//Iterating through the query result and scanning each row into a Post struct
	for res.Next() {
		var post Post
		err = res.Scan(&post.Id, &post.Title, &post.Anons, &post.Full_Text)
		if err != nil {
			// Handling error while scanning database rows
			http.Error(w, "Internal server eror", http.StatusInternalServerError)
			log.Printf("Error scanning database rows: %v", err)
			return
		}
		//Appending the scanned Post struct to the post slice
		posts = append(posts, post)
	}

	//Executing the post template and writing the output to the response writer
	err = t.ExecuteTemplate(w, "post", posts)
	if err != nil {
		// Handling template execution error by returning an internal server error response
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Printf("Error executing template: %v", err)
	}
}

// showPost is an HTTP handler function for displaying a specific article by its ID
// It retrieves the article ID from the request parameters, queries the database for the article
// And renders the article using the show template
func showPost(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var err error
	//Extracing variables from the request parameters
	vars := mux.Vars(r)
	// Parsing HTML template files for the show page, header, and footer
	t, err := template.ParseFiles("templates/show.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		// Handling template parsing error by returning an internal server error response
		http.Error(w, "Error interesting data", http.StatusInternalServerError)
		log.Printf("Error scanning data into database: %v", err)
		return
	}

	//Default page and page size values
	page := 1
	pageSize := 10

	//Checking if a	requested and updating the page variable
	if pageParam := r.FormValue("page"); pageParam != "" {
		page, _ = strconv.Atoi(pageParam)
	}

	//Calculating the offset for pagination
	offset := (page - 1) * pageSize

	//Querying the database for the specific articles using its ID
	res := db.QueryRow("SELECT * FROM articles WHERE id = $1 LIMIT $2 OFFSET $3", vars["id"], pageSize, offset)
	if err != nil {
		// Handling database query error by returning an internal server error response
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Error: %v", err)
		return
	}

	//Creating a Post instance to store the retrieved article data
	showItems = Post{}
	//Scanning the database query result into the showItems variable
	err = res.Scan(&showItems.Id, &showItems.Title, &showItems.Anons, &showItems.Full_Text)
	if err != nil {
		// Handling case when the article is not found
		if err == sql.ErrNoRows {
			http.Error(w, "Article not found", http.StatusNotFound)
			log.Printf("Article not found: %v", err)
		}
		// Handling other errors by returning an internal server error response
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Error: %v", err)
		return
	}
	t.ExecuteTemplate(w, "show", showItems)
}

// complete is an HTTP handler function for displaying a completion page with user-provided data.
// It retrieves the user's name and email from the request form, creates a data map, and renders
// the completion page using the complete template.
func complete(w http.ResponseWriter, r *http.Request) {
	// Parsing HTML template files for the complete page, header, and footer
	t, err := template.ParseFiles("templates/complete.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		// Handling template parsing error by returning an internal server error response
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Printf("Error parsing template files: %v", err)
		return
	}

	//Retrieving user-provided name and email from the request form
	name := r.FormValue("name")
	email := r.FormValue("email")

	//Creating a data map with the name and email
	data := map[string]string{"Name": name, "Email": email}

	// Executing the complete template and writing the output to the response writer
	t.ExecuteTemplate(w, "complete", data)
}

// chat is an HTTP handler function for serving the chat page
// It parses the template files handles potential errors, and executes the chat templates
func chat(w http.ResponseWriter, r *http.Request) {
	//Parsing HTML template files for the chat page, header, and footer
	t, err := template.ParseFiles("templates/chat.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		//Handling template parsing error by returning an internal server error response
		http.Error(w, "Error interesting data", http.StatusInternalServerError)
		log.Printf("Error scanning data into database: %v", err)
		return
	}

	//Executing the chat template and writing the output to the response writer
	t.ExecuteTemplate(w, "chat", nil)
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
		t, err := template.ParseFiles("templates/googleSignIn.html", "templates/header.html", "templates/footer.html")
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
			t, err := template.ParseFiles("templates/complete.html", "templates/header.html", "templates/footer.html")
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

// main function is an the entery point of the application
// It initializes the database connection, sets up third-party authentication providers
// And starts the HTTP server to handle incoming requests
func main() {
	err := godotenv.Load("C:/Arsen/VSC/Projects/VoAr/st.env")
	if err != nil {
		fmt.Println("Error loading .env:", err)
		return
	}
	//Initializing the database connection
	db, err := initDB()
	if err != nil {
		log.Fatal("Error initializing database:", err)
	}
	defer db.Close() //Closing the database connection then main function exits

	//Setting up the Google authentication provider
	google.Google()

	//Creating an HTTP server with the 	configured database connection
	app.Program()

	//Creating the HTTP server with the configured database connection
	server := HandleFunc(db)

	//Starting the HTTP server and handling any potential errors
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Server error:", err)
	}
}
