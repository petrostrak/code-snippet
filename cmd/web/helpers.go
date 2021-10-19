package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/justinas/nosurf"
	"github.com/petrostrak/code-snippet/pkg/models"
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

func (a *application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {

	//Retrive the appropriate template set from the cache based on the page name
	// (like 'home.page.tmpl'). If no entry exists in the cache with the provided
	// name, we call the serverError.
	ts, ok := a.templateCache[name]
	if !ok {
		a.serverError(w, fmt.Errorf("the template %s does not exist", name))
		return
	}

	// Initialize a new buffer
	buf := new(bytes.Buffer)

	// Execute the template set, passing in any dynamic data with
	// the current year injected.
	// Write the template to the buffer, instead of straight to
	// the http.ResponseWriter. If there's an error, call our
	// serverError helper and return
	if err := ts.Execute(buf, a.addDefaultData(td, r)); err != nil {
		a.serverError(w, err)
		return
	}

	// Write the contents of the buffer to the http.ResponseWriter
	buf.WriteTo(w)
}

// Create an addDefaultData helper. This adds the current year to the
// CurrentYear field, and returns the pointer.
func (a *application) addDefaultData(td *templateData, r *http.Request) *templateData {
	if td == nil {
		td = &templateData{}
	}

	// Add the CSRF token to the templateData struct.
	td.CSRFToken = nosurf.Token(r)

	// update our addDefaultData() helper method so that the
	// user ID is automatically added to the templateData struct
	// every time we render a template
	td.AuthenticatedUser = a.authenticatedUser(r)
	td.CurrentYear = time.Now().Year()

	// Add the flash message to the template data, if one exists.
	td.Flash = a.session.PopString(r, "flash")
	return td
}

// The authenticatedUser method returns the ID of the current user from the
// session, or zero if the request is from an unauthenticated user.
func (a *application) authenticatedUser(r *http.Request) *models.User {
	user, ok := r.Context().Value(contextKeyUser).(*models.User)
	if !ok {
		return nil
	}
	return user
}
