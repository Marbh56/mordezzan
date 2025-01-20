package server

import (
	"context"
	"log"
	"net/http"
	"time"

	"database/sql"

	"github.com/marbh56/mordezzan/internal/db"
)

type contextKey string

const UserContextKey contextKey = "user"

func (s *Server) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Public paths that don't require authentication
		publicPaths := map[string]bool{
			"/login":    true,
			"/register": true,
			"/static":   true,
		}

		// Skip authentication for public paths
		if publicPaths[r.URL.Path] {
			next.ServeHTTP(w, r)
			return
		}

		// Get session cookie
		cookie, err := r.Cookie("session_token")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Validate session in database
		queries := db.New(s.db)
		session, err := queries.GetSession(r.Context(), cookie.Value)
		if err != nil {
			if err == sql.ErrNoRows {
				// Invalid or expired session - clear cookie and redirect to login
				http.SetCookie(w, &http.Cookie{
					Name:     "session_token",
					Value:    "",
					Path:     "/",
					HttpOnly: true,
					Expires:  time.Unix(0, 0),
				})
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}
			log.Printf("Database error in auth middleware: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Add user info to request context
		// Create a pointer to the session before storing in context
		ctx := context.WithValue(r.Context(), UserContextKey, &session)

		// Call the next handler with our modified context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Helper function to get user from context
func GetUserFromContext(ctx context.Context) (*db.GetSessionRow, bool) {
	user, ok := ctx.Value(UserContextKey).(*db.GetSessionRow)
	return user, ok
}
