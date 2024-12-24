package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("GET /static", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /character/view/{id}", app.characterView)
	mux.HandleFunc("GET /character/create", app.characterCreate)
	mux.HandleFunc("POST /character/create", app.characterCreatePost)

	return mux
}
