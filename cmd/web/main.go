package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/marbh56/mordezzan/internal/database"
)

func main() {

	db, err := database.OpenDB("./mordezzan.db")
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}
	defer db.Close()
	log.Println("Successfully connected to database!")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello this will be a character manager!")
	})

	fmt.Println("Server starting on :8080...")
	http.ListenAndServe(":8080", nil)
}
