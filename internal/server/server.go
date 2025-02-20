package server

import (
	"database/sql"
	"net/http"
)

// Server represents our HTTP server and its dependencies
type Server struct {
	db *sql.DB
}

// NewServer creates a new instance of Server with the given dependencies
func NewServer(db *sql.DB) *Server {
	return &Server{
		db: db,
	}
}

// Routes sets up and returns all the routes for the server
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
	mux.Handle("/characters/inventory/add", s.AuthMiddleware(http.HandlerFunc(s.HandleAddInventoryItem)))
	mux.Handle("/characters/inventory/remove", s.AuthMiddleware(http.HandlerFunc(s.HandleRemoveInventoryItem)))
	mux.Handle("/characters/edit", s.AuthMiddleware(http.HandlerFunc(s.HandleCharacterEdit)))
	mux.Handle("/characters/xp/update", s.AuthMiddleware(http.HandlerFunc(s.HandleUpdateXP)))
	mux.Handle("/characters/hp/update", s.AuthMiddleware(http.HandlerFunc(s.HandleUpdateHP)))
	mux.Handle("/characters/maxhp/update", s.AuthMiddleware(http.HandlerFunc(s.HandleUpdateMaxHP)))
	mux.Handle("/characters/rest", s.AuthMiddleware(http.HandlerFunc(s.HandleRest)))
	mux.Handle("/characters/currency/update", s.AuthMiddleware(http.HandlerFunc(s.HandleCurrencyUpdate)))

	// Weapon mastery route
	mux.Handle("/characters/masteries", s.AuthMiddleware(http.HandlerFunc(s.HandleWeaponMastery)))

	// Home route - protected by auth middleware
	mux.Handle("/", s.AuthMiddleware(http.HandlerFunc(s.HandleHome)))

	return mux
}
