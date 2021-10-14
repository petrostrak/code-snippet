package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

// The serverError helper writes an error message and stack trace to the errorLog
// then sends a generic 500 Internal Server Error response to the user.
func (a *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	a.errorLog.Output(2, trace)

	// http.StatusText() automatically generates a human-friendly text representation of a
	// given HTTP status code, eg. http.StatusText(400) will give a string "Bad Request".
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// The clientError helper sends a specific status code and corresponding description
// to the user.
func (a *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// For consistency we also implement a notFound helper. This is simply a convenience
// wrapper around clientError which sends a 404 Not Found response to the user.
func (a *application) notFound(w http.ResponseWriter) {
	a.clientError(w, http.StatusNotFound)
}
