package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"
)

// write error and stack trace to errLog and send 500 to user
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errLog.Output(2, trace)  // don't log this line
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// send statuscode and corresponding description (for error) to user
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// a simple notFound helper (404)
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

// adds various default (dynamic) data to templateData
func (app *application) addDefaultData(td *templateData, r *http.Request) *templateData {
	if td == nil {
		td = &templateData{}
	}
	td.Flash = app.session.PopString(r, "flash")
	td.CurrentYear = time.Now().Year()

	return td
}

func (app *application) renderTemplate(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
	ts, ok := app.cache[name]
	if !ok {
		app.serverError(w, fmt.Errorf("Template %s doesn't exist", name))
		return
	}

	// writing template to response in two stage, so that half cooked template
	// (with error) is not shown to user
	buf := new(bytes.Buffer)
	if err := ts.Execute(buf, app.addDefaultData(td, r)); err != nil {
		app.serverError(w, err)
		return
	}

	buf.WriteTo(w)
}
