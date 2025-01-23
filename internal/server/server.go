package server

import (
	"database/sql"
	"net/http"
)

type Server struct {
	db *sql.DB
}

func NewServer(db *sql.DB) *Server {
	return &Server{
		db: db,
	}
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()

	// Static file serving
	fileServer := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	// Auth routes
	mux.HandleFunc("/login", s.HandleLogin)
	mux.HandleFunc("/register", s.HandleRegister)
	mux.HandleFunc("/logout", s.HandleLogout)

	// Character routes - all protected by auth middleware
	mux.Handle("/characters", s.AuthMiddleware(http.HandlerFunc(s.HandleCharacterList)))
	mux.Handle("/characters/create", s.AuthMiddleware(http.HandlerFunc(s.HandleCharacterCreate)))
	mux.Handle("/characters/detail", s.AuthMiddleware(http.HandlerFunc(s.HandleCharacterDetail)))
	mux.Handle("/characters/inventory/add", s.AuthMiddleware(http.HandlerFunc(s.HandleAddInventoryItem))) // Home route - protected by auth middleware
	mux.Handle("/characters/edit", s.AuthMiddleware(http.HandlerFunc(s.HandleCharacterEdit)))
	mux.Handle("/", s.AuthMiddleware(http.HandlerFunc(s.HandleHome)))

	return mux
}
