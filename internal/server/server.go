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

	// Static assets
	fileServer := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	// Authentication routes (public)
	mux.HandleFunc("/login", s.HandleLogin)
	mux.HandleFunc("/register", s.HandleRegister)
	mux.HandleFunc("/logout", s.HandleLogout)

	// Character management routes (protected)
	mux.Handle("/characters", s.AuthMiddleware(http.HandlerFunc(s.HandleCharacterList)))
	mux.Handle("/characters/create", s.AuthMiddleware(http.HandlerFunc(s.HandleCharacterCreate)))
	mux.Handle("/characters/detail", s.AuthMiddleware(http.HandlerFunc(s.HandleCharacterDetail)))
	mux.Handle("/characters/edit", s.AuthMiddleware(http.HandlerFunc(s.HandleCharacterEdit)))
	mux.Handle("/characters/delete", s.AuthMiddleware(http.HandlerFunc(s.HandleDeleteCharacter)))

	// Character status routes (protected)
	mux.Handle("/characters/hp/update", s.AuthMiddleware(http.HandlerFunc(s.HandleUpdateHP)))
	mux.Handle("/characters/hp/form", s.AuthMiddleware(http.HandlerFunc(s.HandleHPForm)))
	mux.Handle("/characters/hp/cancel", s.AuthMiddleware(http.HandlerFunc(s.HandleHPCancel)))
	mux.Handle("/characters/maxhp/update", s.AuthMiddleware(http.HandlerFunc(s.HandleUpdateMaxHP)))
	mux.Handle("/characters/maxhp/form", s.AuthMiddleware(http.HandlerFunc(s.HandleMaxHPForm)))
	mux.Handle("/characters/rest", s.AuthMiddleware(http.HandlerFunc(s.HandleRest)))

	// Currency routes (protected)
	mux.Handle("/characters/currency/update", s.AuthMiddleware(http.HandlerFunc(s.HandleCurrencyUpdate)))

	// Inventory management routes (protected)
	mux.Handle("/characters/inventory/add", s.AuthMiddleware(http.HandlerFunc(s.HandleAddInventoryItem)))
	mux.Handle("/characters/inventory/remove", s.AuthMiddleware(http.HandlerFunc(s.HandleRemoveInventoryItem)))
	mux.Handle("/characters/inventory/update", s.AuthMiddleware(http.HandlerFunc(s.HandleUpdateInventoryItem)))
	mux.Handle("/characters/inventory/equip", s.AuthMiddleware(http.HandlerFunc(s.HandleEquipItem)))
	mux.Handle("/characters/inventory/unequip", s.AuthMiddleware(http.HandlerFunc(s.HandleUnequipItem)))
	mux.Handle("/characters/inventory/move", s.AuthMiddleware(http.HandlerFunc(s.HandleMoveToContainer)))

	// New modal inventory routes
	mux.Handle("/characters/inventory/modal", s.AuthMiddleware(http.HandlerFunc(s.HandleInventoryModal)))
	mux.Handle("/characters/inventory/add-modal", s.AuthMiddleware(http.HandlerFunc(s.HandleAddItemModal)))

	// User settings routes (protected)
	mux.Handle("/settings", s.AuthMiddleware(http.HandlerFunc(s.HandleSettings)))
	mux.Handle("/settings/update", s.AuthMiddleware(http.HandlerFunc(s.HandleUpdateUser)))
	mux.Handle("/settings/update-password", s.AuthMiddleware(http.HandlerFunc(s.HandleUpdatePassword)))
	mux.Handle("/account/delete", s.AuthMiddleware(http.HandlerFunc(s.HandleAccountDelete)))

	// Home page (protected)
	mux.Handle("/", s.AuthMiddleware(http.HandlerFunc(s.HandleHome)))

	return mux
}
