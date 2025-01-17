package server

import (
	"database/sql"
	"log"
	"net/http"
	"text/template"
	"time"
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

	// --------
	mux.HandleFunc("/login", s.HandleLogin)
	//	mux.HandleFunc("/logout", s.HandleLogout)
	//	mux.HandleFunc("/register", s.HandleRegister)

	// --------
	return mux
}

func (s *Server) HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tmpl, err := template.ParseFiles(
		"templates/layout/base.html",
		"templates/auth/login.html",
	)
	if err != nil {
		log.Printf("Template parsing error: %v", err)
		http.Error(w, "Internal Server Error", 500)
		return
	}

	data := struct {
		IsAuthenticated bool
		FlashMessage    string
		CurrentYear     int
	}{
		IsAuthenticated: false,
		FlashMessage:    "",
		CurrentYear:     time.Now().Year(),
	}

	err = tmpl.ExecuteTemplate(w, "base.html", data)
	if err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "Internal Server Error", 500)
		return
	}
}
