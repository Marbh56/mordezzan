package main

import (
	"fmt"
	"net/http"

	"github.com/marbh56/mordezzan/internal/db"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	tmpl, ok := app.templateCache["home.tmpl"]
	if !ok {
		app.serverError(w, fmt.Errorf("the template %s does not exist", "home.tmpl"))
		return
	}

	err := tmpl.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) registerForm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	tmpl, ok := app.templateCache["registration.tmpl"]
	if !ok {
		app.serverError(w, fmt.Errorf("the template %s does not exist", "register.tmpl"))
		return
	}

	err := tmpl.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	_, err = app.models.Users.CreateUser(r.Context(), db.CreateUserParams{
		Name:         name,
		Email:        email,
		PasswordHash: hashedPassword,
	})

	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
