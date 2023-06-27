package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	// "text/template"

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

	app.renderTemplate(w, r, "show.page.tmpl", &templateData{
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

	errMap := make(map[string]string)

	// validation for title
	if strings.TrimSpace(title) == "" {
		errMap["title"] = "Title can't be empty"
	} else if utf8.RuneCountInString(title) > 100 {
		errMap["title"] = "Title length can't be more than 100 characters"
	}

	// validation for content
	if strings.TrimSpace("content") == "" {
		errMap["content"] = "Content can't be empty"
	}

	// validation for expiry (just to be sure)
	if strings.TrimSpace(expires) == "" {
		errMap["expires"] = "Expiry duration can't be empty"
	}

	if len(errMap) > 0 {
		fmt.Fprint(w, errMap)
		return
	}

	id, err := app.gomi.Create(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// redirect to see the gomi
	http.Redirect(w, r, fmt.Sprintf("/gomi?id=%d", id), http.StatusSeeOther)
}

func (app *application) createGomiForm(w http.ResponseWriter, r *http.Request) {
	app.renderTemplate(w, r, "create.page.tmpl", nil)
}
