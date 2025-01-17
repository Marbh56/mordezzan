package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/marbh56/mordezzan/internal/database"
	"github.com/marbh56/mordezzan/internal/server"
)

func main() {

	db, err := database.OpenDB("./mordezzan.db")
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}
	defer db.Close()
	log.Println("Successfully connected to database!")

	server := server.NewServer(db)

	fmt.Println("Server starting on :8080...")
	http.ListenAndServe(":8080", server.Routes())
}
