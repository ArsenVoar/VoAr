package app

import (
	"database/sql"
	"errors"
	"html/template"
	"log"
	"net/http"

	"github.com/lib/pq"
	"github.com/markbates/goth"
)

func SaveUsersToDB(w http.ResponseWriter, r *http.Request, db *sql.DB, user goth.User) error {
	_, err := db.Exec("INSERT INTO users (name, email) VALUES ($1, $2)", user.Name, user.Email)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" { // PostgreSQL error
				return errors.New("user with the same name or email already exists")
			}
		}
		// Other error
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Printf("Error saving user data to the database: %v", err)
		return err
	}
	return nil
}

// userSavedSuccesfull handles the case where user data is succesfully saved to the database
// It parses the template files for the userSavedSuccesfull, header, and footer,
// Handles potential errors, and executes the create page template.
func UserSavedSuccesfull(w http.ResponseWriter, r *http.Request) {
	// Parsing HTML template files for the userSavedSuccesfull, header, and footer
	t, err := template.ParseFiles("web/templates/userSavedSuccesfull.html", "web//templates/header.html", "web/templates/footer.html")
	if err != nil {
		// Handling template parsing error by returning an internal server error response
		http.Error(w, "Error interesting data", http.StatusInternalServerError)
		log.Printf("Error scanning data into database: %v", err)
		return
	}

	// Executing the userSavedSuccesfull template and writing the output to the response writer
	t.ExecuteTemplate(w, "userSavedSuccesfull", nil)
}

// userExists handles the case where a user with the same name or email already exists
// It parses the template files for the userExists, header, and footer
// Handles potential errors, and executes the create page template
func UserExists(w http.ResponseWriter, r *http.Request) {
	// Parsing HTML template files for the userExists, header, and footer
	t, err := template.ParseFiles("web/templates/userExists.html", "web/templates/header.html", "web/templates/footer.html")
	if err != nil {
		// Handling template parsing error by returning an internal server error response
		http.Error(w, "Error interesting data", http.StatusInternalServerError)
		log.Printf("Error scanning data into database: %v", err)
		return
	}

	// Executing the userExists page template and writing the output to the response writer
	t.ExecuteTemplate(w, "userExists", nil)
}

// complete is an HTTP handler function for displaying a completion page with user-provided data.
// It retrieves the user's name and email from the request form, creates a data map, and renders
// the completion page using the complete template.
func Complete(w http.ResponseWriter, r *http.Request) {
	// Parsing HTML template files for the complete page, header, and footer
	t, err := template.ParseFiles("web/templates/complete.html", "web/templates/header.html", "web/templates/footer.html")
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
