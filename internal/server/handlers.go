package server

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/marbh56/mordezzan/internal/db"
	"golang.org/x/crypto/bcrypt"
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
	mux.HandleFunc("/register", s.HandleRegister)

	// --------
	return mux
}

func (s *Server) HandleRegister(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.handleRegisterForm(w, r)
	case http.MethodPost:
		s.handleRegistrerSubmission(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleRegisterForm(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		"templates/layout/base.html",
		"templates/auth/registration.html",
	)

	if err != nil {
		log.Printf("Template parsing error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	data := struct {
		IsAuthenticated bool
		FlashMessage    string
		CurrentYear     int
	}{
		IsAuthenticated: false,
		FlashMessage:    r.URL.Query().Get("error"),
		CurrentYear:     time.Now().Year(),
	}

	err = tmpl.ExecuteTemplate(w, "base.html", data)
	if err != nil {
		log.Printf("Tempalate execution error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleRegistrerSubmission(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	username := r.Form.Get("username")
	email := r.Form.Get("email")
	password := r.Form.Get("password")
	confirmPassword := r.Form.Get("confirm-password")

	if username == "" || email == "" || password == "" || confirmPassword == "" {
		http.Redirect(w, r, "/register?error=All fields are required", http.StatusSeeOther)
		return
	}

	if len(password) < 8 {
		http.Redirect(w, r, "/register?error=Password must be at least 8 characters", http.StatusSeeOther)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	queries := db.New(s.db)

	params := db.CreateUserParams{
		Username:     username,
		Email:        email,
		PasswordHash: string(hashedPassword),
	}

	_, err = queries.CreateUser(context.Background(), params)
	if err != nil {
		log.Printf("Error creating user: %v", err)

		if err == sql.ErrNoRows {
			http.Redirect(w, r, "/register?error=Username or email already exists", http.StatusSeeOther)
		}

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/login?message=Registration successful! Please log in", http.StatusSeeOther)
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
