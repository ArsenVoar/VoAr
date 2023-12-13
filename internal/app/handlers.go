package app

import (
	"html/template"
	"log"
	"net/http"
)

// contextKey is a custom type for convience when using it in yhe context
type ContextKey string

// dbkey is the context key for passing the database instance to handlers
const DbKey ContextKey = "db"

// User represents a structure for storing user data
type User struct {
	ID    int
	Name  string
	Email string
}

// Post represents a structure for storing arcticle data
type Pst struct {
	Id        int    //Unique identifier for the arcticle
	Title     string //Title of the arcticle
	Anons     string //Brief summary or announcement of the article
	Full_Text string //Full text content of the article
}

var (
	posts     = []Pst{} //Global slice to store multiple article instances
	showItems = Pst{}   //Global varriable to store a single article instance for display
)

// mainPage is an HTTP handler function for serving the main page.
// It parses the template files for the main page, header, and footer,
// Handles potential errors, and executes the main page template
func MainPage(w http.ResponseWriter, r *http.Request) {
	// Parsing HTML template files for the main page, header, and footer
	// Замените строки загрузки шаблонов в mainPage, chat и complete функциях
	t, err := template.ParseFiles("web/templates/mainPage.html", "web/templates/header.html", "web/templates/footer.html")
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
func Examples(w http.ResponseWriter, r *http.Request) {
	// Parsing HTML template files for the examples page, header, and footer
	t, err := template.ParseFiles("web/templates/examples.html", "web/templates/header.html", "web/templates/footer.html")
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
func Create(w http.ResponseWriter, r *http.Request) {
	// Parsing HTML template files for the create page, header, and footer
	t, err := template.ParseFiles("web/templates/create.html", "web/templates/header.html", "web/templates/footer.html")
	if err != nil {
		// Handling template parsing error by returning an internal server error response
		http.Error(w, "Error interesting data", http.StatusInternalServerError)
		log.Printf("Error scanning data into database: %v", err)
		return
	}

	// Executing the create page template and writing the output to the response writer
	t.ExecuteTemplate(w, "create", nil)
}

// chat is an HTTP handler function for serving the chat page
// It parses the template files handles potential errors, and executes the chat templates
func Chat(w http.ResponseWriter, r *http.Request) {
	//Parsing HTML template files for the chat page, header, and footer
	t, err := template.ParseFiles("web/templates/chat.html", "web/templates/header.html", "web/templates/footer.html")
	if err != nil {
		//Handling template parsing error by returning an internal server error response
		http.Error(w, "Error interesting data", http.StatusInternalServerError)
		log.Printf("Error scanning data into database: %v", err)
		return
	}

	//Executing the chat template and writing the output to the response writer
	t.ExecuteTemplate(w, "chat", nil)
}
