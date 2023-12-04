package app

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// contextKey is a custom type for convience when using it in yhe context
type contextKey string

// dbkey is the context key for passing the database instance to handlers
const dbKey contextKey = "db"

// Post represents a structure for storing arcticle data
type Post struct {
	Id        int    //Unique identifier for the arcticle
	Title     string //Title of the arcticle
	Anons     string //Brief summary or announcement of the article
	Full_Text string //Full text content of the article
}

var (
	posts     = []Post{} //Global slice to store multiple article instances
	showItems = Post{}   //Global varriable to store a single article instance for display
)

// mainPage is an HTTP handler function for serving the main page.
// It parses the template files for the main page, header, and footer,
// Handles potential errors, and executes the main page template
func mainPage(w http.ResponseWriter, r *http.Request) {
	// Parsing HTML template files for the main page, header, and footer
	// Замените строки загрузки шаблонов в mainPage, chat и complete функциях
	t, err := template.ParseFiles("cmd/pkg/app/templates/mainPage.html", "cmd/pkg/app/templates/header.html", "cmd/pkg/app/templates/footer.html")
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
	t, err := template.ParseFiles("cmd/pkg/app/templates/examples.html", "cmd/pkg/app/templates/header.html", "cmd/pkg/app/templates/footer.html")
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
	t, err := template.ParseFiles("cmd/pkg/app/templates/create.html", "cmd/pkg/app/templates/header.html", "cmd/pkg/app/templates/footer.html")
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
	t, err := template.ParseFiles("cmd/pkg/app/templates/post.html", "cmd/pkg/app/templates/show.html", "cmd/pkg/app/templates/header.html", "cmd/pkg/app/templates/footer.html")
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
	t, err := template.ParseFiles("cmd/pkg/app/templates/show.html", "cmd/pkg/app/templates/header.html", "cmd/pkg/app/templates/footer.html")
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
	t, err := template.ParseFiles("cmd/pkg/app/templates/complete.html", "cmd/pkg/app/templates/header.html", "cmd/pkg/app/templates/footer.html")
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
	t, err := template.ParseFiles("cmd/pkg/app/templates/chat.html", "cmd/pkg/app/templates/header.html", "cmd/pkg/app/templates/footer.html")
	if err != nil {
		//Handling template parsing error by returning an internal server error response
		http.Error(w, "Error interesting data", http.StatusInternalServerError)
		log.Printf("Error scanning data into database: %v", err)
		return
	}

	//Executing the chat template and writing the output to the response writer
	t.ExecuteTemplate(w, "chat", nil)
}
