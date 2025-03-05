package server

import (
	"net/http"
	"time"
)

func (s *Server) HandleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	var username string
	isAuthenticated := false

	if user, ok := GetUserFromContext(r.Context()); ok {
		username = user.Username
		isAuthenticated = true
	}

	data := struct {
		IsAuthenticated bool
		Username        string
		FlashMessage    string
		CurrentYear     int
	}{
		IsAuthenticated: isAuthenticated,
		Username:        username,
		FlashMessage:    r.URL.Query().Get("message"),
		CurrentYear:     time.Now().Year(),
	}
	RenderTemplate(w, "templates/home.html", "base.html", data)
}
