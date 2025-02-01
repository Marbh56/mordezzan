package main

import (
	"log"
	"net/http"

	"github.com/marbh56/mordezzan/internal/database"
	"github.com/marbh56/mordezzan/internal/server"
)

func main() {
	log.Println("Starting application...")

	log.Println("Attempting to open database...")
	db, err := database.OpenDB("./mordezzan.db")
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}
	defer db.Close()
	log.Println("Successfully connected to database!")

	log.Println("Creating new server instance...")
	srv := server.NewServer(db)

	log.Println("Setting up routes...")
	handler := srv.Routes()

	log.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
