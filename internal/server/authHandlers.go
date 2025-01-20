package server

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/marbh56/mordezzan/internal/db"
	"golang.org/x/crypto/bcrypt"
)

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

	if password != confirmPassword {
		http.Redirect(w, r, "/register?error=Passwords do not match", http.StatusSeeOther)
		return
	}

	if len(password) < 8 {
		http.Redirect(w, r, "/register?error=Password must be at least 8 characters", http.StatusSeeOther)
		return
	}

	queries := db.New(s.db)

	_, err := queries.GetUserByUsername(context.Background(), username)
	if err != nil {
		http.Redirect(w, r, "/register?error=Username already exists", http.StatusSeeOther)
		return
	} else if err != sql.ErrNoRows {
		log.Printf("Error checking username: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	_, err = queries.GetUserByEmail(context.Background(), email)
	if err == nil {
		http.Redirect(w, r, "register?error=Email already exists", http.StatusSeeOther)
		return
	} else if err != sql.ErrNoRows {
		log.Printf("Error checking email: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

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
	switch r.Method {
	case http.MethodGet:
		s.handleLoginForm(w, r)
	case http.MethodPost:
		s.handleLoginSubmission(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleLoginForm(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		"templates/layout/base.html",
		"templates/auth/login.html",
	)
	if err != nil {
		log.Printf("Template parsing error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := struct {
		IsAuthenticated bool
		FlashMessage    string
		CurrentYear     int
	}{
		IsAuthenticated: false,
		FlashMessage:    r.URL.Query().Get("message"),
		CurrentYear:     time.Now().Year(),
	}

	err = tmpl.ExecuteTemplate(w, "base.html", data)
	if err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleLoginSubmission(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	username := r.Form.Get("username")
	password := r.Form.Get("password")

	if username == "" || password == "" {
		http.Redirect(w, r, "/login?message=Username and password are required", http.StatusSeeOther)
		return
	}

	queries := db.New(s.db)

	// Get user from database
	user, err := queries.GetUserByUsername(r.Context(), username)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Redirect(w, r, "/login?message=Invalid username or password", http.StatusSeeOther)
			return
		}
		log.Printf("Database error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		http.Redirect(w, r, "/login?message=Invalid username or password", http.StatusSeeOther)
		return
	}

	// Generate session token
	token := make([]byte, 32)
	_, err = rand.Read(token)
	if err != nil {
		log.Printf("Error generating session token: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	sessionToken := hex.EncodeToString(token)

	// Create session in database
	sessionParams := db.CreateSessionParams{
		Token:     sessionToken,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(24 * time.Hour), // Sessions last 24 hours
	}

	_, err = queries.CreateSession(r.Context(), sessionParams)
	if err != nil {
		log.Printf("Error creating session: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Set session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   r.TLS != nil, // Set Secure flag if using HTTPS
		SameSite: http.SameSiteLaxMode,
		Expires:  sessionParams.ExpiresAt,
	})

	// Redirect to home page after successful login
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *Server) HandleLogout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie("session_token")
	if err == nil {
		queries := db.New(s.db)
		err = queries.DeleteSession(r.Context(), cookie.Value)
		if err != nil {
			log.Printf("Error deleting session: %v", err)
		}
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Unix(0, 0),
	})

	http.Redirect(w, r, "/login?message=Successfully logged out", http.StatusSeeOther)
}
