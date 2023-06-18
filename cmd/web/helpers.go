package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
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
