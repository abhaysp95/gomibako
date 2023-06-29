package main

import (
	"fmt"
	"net/http"
	"strconv"

	// "text/template"

	"github.com/abhaysp95/gomibako/pkg/forms"
	"github.com/abhaysp95/gomibako/pkg/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
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
		app.errLog.Println(err)
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

	flash := app.session.PopString(r, "flash")

	app.renderTemplate(w, r, "show.page.tmpl", &templateData{
		Flash: flash,
		Gomi: g,
	})
}

// handler to create new gomi
func (app *application) createGomi(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")
	expires := r.PostForm.Get("expires")

	form := forms.New(r.PostForm)

	// validation for required
	form.Required("title", "content")
	form.MaxLength("title", 100)
	form.PermittedValues("expires", "1", "7", "365")

	if !form.Valid() {
		app.errLog.Println("validation error", form.ErrMap)
		app.renderTemplate(w, r, "create.page.tmpl", &templateData{
			Form: form,
		})
		return
	}

	id, err := app.gomi.Create(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// add key-value pair to session data, if there's no existing sesion for
	// current user then a new session, empty session will be created by
	// session middleware
	app.session.Put(r, "flash", "Gomi created successfully")

	// redirect to see the gomi
	http.Redirect(w, r, fmt.Sprintf("/gomi?id=%d", id), http.StatusSeeOther)
}

func (app *application) createGomiForm(w http.ResponseWriter, r *http.Request) {
	form := forms.New(nil)
	app.renderTemplate(w, r, "create.page.tmpl", &templateData{
		Form: form,
	})
}
