package main

import (
	"VoAr/internal/database"
	"VoAr/internal/handler"
	google "VoAr/pkg/google"
	"log"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load("st.env")
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	db, err := database.InitDB()
	if err != nil {
		log.Fatal("Error initializing database:", err)
	}
	defer db.Close()

	google.Google()

	server := handler.HandleFunc(db)

	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Server error:", err)
	}
}
