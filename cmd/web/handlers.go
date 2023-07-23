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

// handler to sign up new user
func (app *application) signupUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.serverError(w, err)
		return
	}

	name := r.PostForm.Get("name")
	email := r.PostForm.Get("email")
	passwd := r.PostForm.Get("passwd")

	form := forms.New(r.PostForm)
	form.Required("name")
	form.Required("email")
	form.Required("passwd")
	form.MinLength("passwd", 10)
	form.MatchesPattern("email", forms.EmailRx)

	if !form.Valid() {
		app.errLog.Println("validation error", form.ErrMap)
		app.renderTemplate(w, r, "signup.page.tmpl", &templateData{
			Form: form,
		})
		return
	}

	err = app.user.Insert(name, email, passwd)
	if err == models.ErrDuplicateEmail {
		form.ErrMap.Add("email", err.Error())
		app.renderTemplate(w, r, "signup.page.tmpl", &templateData{
			Form: form,
		})
	} else if err == models.ErrInvalidCredentials {
		form.ErrMap.Add("passwd", err.Error())
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	app.session.Put(r, "flash", "Account created successfully. Please log in")

	// redirect to the home page
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

// handler to show sign up form
func (app *application) signupUserForm(w http.ResponseWriter, r *http.Request) {
	form := forms.New(nil)
	app.renderTemplate(w, r, "signup.page.tmpl", &templateData{
		Form: form,
	})
}

// handler to login user
func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)

	id, err := app.user.Authenticate(form.Get("email"), form.Get("passwd"))
	if err == models.ErrInvalidCredentials {
		form.ErrMap.Add("generic", "Login failed. Email or Password is incorrect")
		app.renderTemplate(w, r, "login.page.tmpl", &templateData{
			Form: form,
		})
		return
	}
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.session.Put(r, "userId", id)

	http.Redirect(w, r, "/gomi/create", http.StatusSeeOther)
}

// handler to show login form
func (app *application) loginUserForm(w http.ResponseWriter, r *http.Request) {
	form := forms.New(nil)
	app.renderTemplate(w, r, "login.page.tmpl", &templateData{
		Form: form,
	})
}

// handler to logout user
func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	app.infoLog.Println("logout user")
}
