package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	r := chi.NewRouter()

	r.Get("/", app.home)
	r.Get("/gomi", app.showGomi)   // different from what
	r.Get("/gomi/create", app.createGomiForm)
	r.Post("/gomi/create", app.createGomi)

	// file server to serve static files
	fileServer := http.FileServer(http.Dir("./ui/static"))

	r.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return app.recoverPanic(app.requestLoggin(secureHeaders(r)))
}
