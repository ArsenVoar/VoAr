package main

import (
	"VoAr/app"               //Importing the app package (assuming it handles application-specific logic)
	google "VoAr/pkg/google" //Importing the Google package for authentication
	"log"                    //Package for logging

	"github.com/joho/godotenv" //Package for loading environment variables from a .env file.
	_ "github.com/lib/pq"      //PostgreSQL driver for the database/sql package
)

// main function is an the entery point of the application
// It initializes the database connection, sets up third-party authentication providers
// And starts the HTTP server to handle incoming requests
func main() {
	// Load environment variables from the specified file
	err := godotenv.Load("st.env")
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	//Initializing the database connection
	db, err := app.InitDB()
	if err != nil {
		log.Fatal("Error initializing database:", err)
	}
	defer db.Close() //Closing the database connection then main function exits

	//Setting up the Google authentication provider
	google.Google()

	//Creating the HTTP server with the configured database connection
	server := app.HandleFunc(db)

	//Starting the HTTP server and handling any potential errors
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Server error:", err)
	}
}
