package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	r := chi.NewRouter()

	r.Get("/", app.session.Enable(http.HandlerFunc(app.home)).(http.HandlerFunc))
	r.Get("/gomi", app.session.Enable(http.HandlerFunc(app.showGomi)).(http.HandlerFunc))   // different from what
	r.Get("/gomi/create", app.session.Enable(http.HandlerFunc(app.createGomiForm)).(http.HandlerFunc))
	r.Post("/gomi/create", app.session.Enable(http.HandlerFunc(app.createGomi)).(http.HandlerFunc))

	// file server to serve static files
	fileServer := http.FileServer(http.Dir("./ui/static"))

	r.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return app.recoverPanic(app.requestLoggin(secureHeaders(r)))
}
