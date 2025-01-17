package main

import (
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			app.registerForm(w, r)
		case http.MethodPost:
			app.register(w, r)
		default:
			w.Header().Set("Allow", http.MethodGet+", "+http.MethodPost)
			app.clientError(w, http.StatusMethodNotAllowed)
		}
	})

	return mux
}
