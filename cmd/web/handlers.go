package main

import (
	"fmt"
	"net/http"
	"strconv"

	// "text/template"

	"github.com/abhaysp95/gomibako/pkg/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	gl, err := app.gomi.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.renderTemplate(w, r, "home.page.tmpl", &templateData{
		GomiList: gl,
	})
}

// handler to showing individual gomi
func (app *application) showGomi(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	g, err := app.gomi.Get(id)
	if err == models.ErrNoRecord {
		app.notFound(w)
		return
	} else if err != nil {
		app.serverError(w, err)
	}

	app.renderTemplate(w, r, "show.page.tmpl", &templateData{
		Gomi: g,
	})
}

// handler to create new gomi
func (app *application) createGomi(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	title := "Title from Gomi appl."
	content := "A small content\nJust for demonstration purpose\n\n - gomibako"
	expires := "5"

	id, err := app.gomi.Create(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// redirect to see the gomi
	http.Redirect(w, r, fmt.Sprintf("/gomi?id=%d", id), http.StatusSeeOther)
}
