package main

import (
	"net"
	"net/http"

	"github.com/marbh56/mordezzan/internal/database"
	"github.com/marbh56/mordezzan/internal/logger"
	"github.com/marbh56/mordezzan/internal/server"
	"go.uber.org/zap"
)

func main() {
	// Initialize logger
	if err := logger.Initialize("development"); err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}
	defer logger.Sync()

	logger.Info("Starting application...")

	// Get the local IP address
	ip, err := getLocalIP()
	if err != nil {
		logger.Fatal("Failed to get local IP address", zap.Error(err))
	}

	logger.Info("Attempting to open database...")
	db, err := database.OpenDB("./mordezzan.db")
	if err != nil {
		logger.Fatal("Failed to open database", zap.Error(err))
	}
	defer db.Close()
	logger.Info("Successfully connected to database")

	logger.Info("Creating new server instance...")
	srv := server.NewServer(db)

	logger.Info("Setting up routes...")
	handler := srv.Routes()

	// Construct the address with local IP
	addr := ip + ":8080"
	logger.Info("Server starting", zap.String("address", addr))

	if err := http.ListenAndServe("localhost:8080", handler); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
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
