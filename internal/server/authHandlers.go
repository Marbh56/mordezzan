package server

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/marbh56/mordezzan/internal/db"
	"github.com/marbh56/mordezzan/internal/logger"
	"go.uber.org/zap"
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
		logger.Error("Failed to parse registration template",
			zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
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
		logger.Error("Failed to execute registration template",
			zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleRegistrerSubmission(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		logger.Error("Failed to parse registration form",
			zap.Error(err))
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	username := r.Form.Get("username")
	email := r.Form.Get("email")
	password := r.Form.Get("password")
	confirmPassword := r.Form.Get("confirm-password")

	// Validate required fields
	if username == "" || email == "" || password == "" || confirmPassword == "" {
		logger.Warn("Registration attempt with missing required fields",
			zap.String("username", username),
			zap.String("email", email))
		http.Redirect(w, r, "/register?error=All fields are required", http.StatusSeeOther)
		return
	}

	// Validate password match
	if password != confirmPassword {
		logger.Warn("Registration password mismatch",
			zap.String("username", username),
			zap.String("email", email))
		http.Redirect(w, r, "/register?error=Passwords do not match", http.StatusSeeOther)
		return
	}

	// Validate password length
	if len(password) < 8 {
		logger.Warn("Registration attempt with short password",
			zap.String("username", username),
			zap.String("email", email),
			zap.Int("password_length", len(password)))
		http.Redirect(w, r, "/register?error=Password must be at least 8 characters", http.StatusSeeOther)
		return
	}

	queries := db.New(s.db)

	// Check if username exists
	_, err := queries.GetUserByUsername(context.Background(), username)
	if err == nil {
		logger.Warn("Registration attempt with existing username",
			zap.String("username", username))
		http.Redirect(w, r, "/register?error=Username already exists", http.StatusSeeOther)
		return
	} else if err != sql.ErrNoRows {
		logger.Error("Database error checking username",
			zap.Error(err),
			zap.String("username", username))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Check if email exists
	_, err = queries.GetUserByEmail(context.Background(), email)
	if err == nil {
		logger.Warn("Registration attempt with existing email",
			zap.String("email", email))
		http.Redirect(w, r, "register?error=Email already exists", http.StatusSeeOther)
		return
	} else if err != sql.ErrNoRows {
		logger.Error("Database error checking email",
			zap.Error(err),
			zap.String("email", email))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("Failed to hash password",
			zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Create user
	params := db.CreateUserParams{
		Username:     username,
		Email:        email,
		PasswordHash: string(hashedPassword),
	}

	_, err = queries.CreateUser(context.Background(), params)
	if err != nil {
		logger.Error("Failed to create user",
			zap.Error(err),
			zap.String("username", username),
			zap.String("email", email))

		if err == sql.ErrNoRows {
			http.Redirect(w, r, "/register?error=Username or email already exists", http.StatusSeeOther)
			return
		}

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	logger.Info("New user registered successfully",
		zap.String("username", username),
		zap.String("email", email))

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
		logger.Error("Failed to parse login template",
			zap.Error(err))
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
		logger.Error("Failed to execute login template",
			zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleLoginSubmission(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		logger.Error("Failed to parse login form",
			zap.Error(err))
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	username := r.Form.Get("username")
	password := r.Form.Get("password")

	if username == "" || password == "" {
		logger.Warn("Login attempt with missing credentials",
			zap.String("username", username))
		http.Redirect(w, r, "/login?message=Username and password are required", http.StatusSeeOther)
		return
	}

	queries := db.New(s.db)

	// Get user from database
	user, err := queries.GetUserByUsername(r.Context(), username)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Warn("Login attempt with non-existent username",
				zap.String("username", username))
			http.Redirect(w, r, "/login?message=Invalid username or password", http.StatusSeeOther)
			return
		}
		logger.Error("Database error during login",
			zap.Error(err),
			zap.String("username", username))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		logger.Warn("Failed login attempt - invalid password",
			zap.String("username", username))
		http.Redirect(w, r, "/login?message=Invalid username or password", http.StatusSeeOther)
		return
	}

	// Generate session token
	token := make([]byte, 32)
	_, err = rand.Read(token)
	if err != nil {
		logger.Error("Failed to generate session token",
			zap.Error(err),
			zap.String("username", username))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	sessionToken := hex.EncodeToString(token)

	// Create session
	sessionParams := db.CreateSessionParams{
		Token:     sessionToken,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	_, err = queries.CreateSession(r.Context(), sessionParams)
	if err != nil {
		logger.Error("Failed to create session",
			zap.Error(err),
			zap.String("username", username),
			zap.Int64("user_id", user.ID))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Set session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   r.TLS != nil,
		SameSite: http.SameSiteLaxMode,
		Expires:  sessionParams.ExpiresAt,
	})

	logger.Info("User logged in successfully",
		zap.String("username", username),
		zap.Int64("user_id", user.ID))

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
			logger.Error("Failed to delete session",
				zap.Error(err),
				zap.String("session_token", cookie.Value))
		}
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Unix(0, 0),
	})

	logger.Info("User logged out")
	http.Redirect(w, r, "/login?message=Successfully logged out", http.StatusSeeOther)
}

func (s *Server) HandleAccountDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user, ok := GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	queries := db.New(s.db)

	// Soft delete user
	err := queries.SoftDeleteUser(r.Context(), user.UserID)
	if err != nil {
		logger.Error("Failed to delete user", zap.Error(err))
		http.Redirect(w, r, "/?message=Error deactivating account", http.StatusSeeOther)
		return
	}

	// Clear session
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Unix(0, 0),
	})

	http.Redirect(w, r, "/?message=Account deactivated successfully", http.StatusSeeOther)
}

func (s *Server) HandleUpdateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user, ok := GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	newEmail := r.Form.Get("email")
	if newEmail == "" {
		http.Redirect(w, r, "/settings?message=Email cannot be empty", http.StatusSeeOther)
		return
	}

	queries := db.New(s.db)
	_, err := queries.UpdateUser(r.Context(), db.UpdateUserParams{
		ID:    user.UserID,
		Email: newEmail,
	})
	if err != nil {
		logger.Error("Failed to update user",
			zap.Error(err),
			zap.String("user_id", strconv.FormatInt(user.UserID, 10)))
		http.Redirect(w, r, "/settings?message=Error updating email", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/settings?message=Email updated successfully", http.StatusSeeOther)
}

func (s *Server) HandleSettings(w http.ResponseWriter, r *http.Request) {
	user, ok := GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	tmpl, err := template.ParseFiles(
		"templates/layout/base.html",
		"templates/auth/settings.html",
	)
	if err != nil {
		logger.Error("Template parsing error", zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := struct {
		IsAuthenticated bool
		Username        string
		User            *db.GetSessionRow
		FlashMessage    string
		CurrentYear     int
	}{
		IsAuthenticated: true,
		Username:        user.Username,
		User:            user,
		FlashMessage:    r.URL.Query().Get("message"),
		CurrentYear:     time.Now().Year(),
	}

	if err := tmpl.ExecuteTemplate(w, "base.html", data); err != nil {
		logger.Error("Template execution error", zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
