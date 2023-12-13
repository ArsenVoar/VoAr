package app

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// save_article is an HTTP handler function for saving an article to the database
// It retrieves form values from the request, validates them, and inserts the data into the database
func Save_article(w http.ResponseWriter, r *http.Request) {
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
	db := r.Context().Value(DbKey).(*sql.DB)

	//Executing the SQL query to insert the article into the database and getting the result
	result, err := db.Exec("INSERT INTO articles (title, anons, full_text, user_id) VALUES ($1, $2, $3) RETURNING id", title, anons, full_text)
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
func Post(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Parsing HTML template files for the post page, show page, header, and footer
	t, err := template.ParseFiles("web/templates/post.html", "web/templates/show.html", "web/templates/header.html", "web/templates/footer.html")
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
	posts = []Pst{}
	//Iterating through the query result and scanning each row into a Post struct
	for res.Next() {
		var post Pst
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
func ShowPost(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var err error
	//Extracing variables from the request parameters
	vars := mux.Vars(r)
	// Parsing HTML template files for the show page, header, and footer
	t, err := template.ParseFiles("web/templates/show.html", "web/templates/header.html", "web/templates/footer.html")
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
	showItems = Pst{}
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
