package main

import (
	"log"
	"net"
	"net/http"

	"github.com/marbh56/mordezzan/internal/database"
	"github.com/marbh56/mordezzan/internal/server"
)

func main() {
	log.Println("Starting application...")

	// Get the local IP address
	ip, err := getLocalIP()
	if err != nil {
		log.Fatal("Failed to get local IP address:", err)
	}

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

	// Construct the address with your local IP
	addr := ip + ":8080"
	log.Printf("Server starting on %s...", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// getLocalIP returns the non loopback local IP of the host
func getLocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback then return it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}
	return "", err
}
