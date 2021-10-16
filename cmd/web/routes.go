package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

var (
	// http.ServeMux is also a handler, which instead of providing
	// a response itself passes the request on to a second handler.
	mux *pat.PatternServeMux
)

func init() {
	mux = pat.New()
}

// Update the signature for the routes() so that it returnst a
// http.Handler instead of a *http.ServeMux
func (a *application) routes() http.Handler {

	// Create a middleware chan containing our 'standard' middleware
	// which will be used for every request our app receives.
	standardMiddleware := alice.New(a.recoverPanic, a.logRequest, secureHeaders)

	// Use the http.NewServeMux() function to initialize a
	// new servemux, then register the home function as the
	// handler for the "/" URL pattern.
	mux.Get("/", http.HandlerFunc(a.home))
	mux.Get("/snippet/create", http.HandlerFunc(a.createSnippetForm))
	mux.Post("/snippet/create", http.HandlerFunc(a.createSnippet))
	mux.Get("/snippet/:id", http.HandlerFunc(a.showSnippet))

	// Create a file server which serves files out of the ./ui/static/ dir.
	// Note that the path given to the http.Dir() function is relative to the
	// project directory root.
	fs := http.FileServer(http.Dir("./ui/static/"))

	// Use the mux.Handle() to register the file server as the handler for
	// all URL paths that start with "/static". When the handler receives a
	// request, it will remove the leading slash from the URL path and then
	// search the ./ui/static directory for the corresponding file to send
	// to the user. So, for this to work correctly, we must strip the leading
	// "/static" from the URL path before passing it to http.FileServer.
	mux.Get("/static/", http.StripPrefix("/static", fs))

	// Return the 'standard' middleware chain followed by the serveMux.
	return standardMiddleware.Then(mux)
}
