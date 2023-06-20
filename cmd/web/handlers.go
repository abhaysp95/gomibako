package main

import (
	"fmt"
	"html/template"
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

	/* gs, err := app.gomi.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}
	for _, gomi := range(gs) {
		fmt.Fprintf(w, "%v\n", gomi)
	} */

	files := []string{
		"./ui/html/home.page.tmpl",
		"./ui/html/base.layout.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.errLog.Println(err.Error())
		app.serverError(w, err)
		return
	}

	err = ts.Execute(w, nil)
	if err != nil {
		app.errLog.Println(err.Error())
		app.serverError(w, err)
	}
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

	files := []string {
		"./ui/html/show.page.tmpl",
		"./ui/html/base.layout.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.errLog.Println(err.Error())
		app.serverError(w, err)
		return
	}

	if err = ts.Execute(w, g); err != nil {
		app.errLog.Println(err.Error())
		app.serverError(w, err)
	}
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
