package main

import "net/http"

var (
	// http.ServeMux is also a handler, which instead of providing
	// a response itself passes the request on to a second handler.
	mux *http.ServeMux
)

func init() {
	mux = http.NewServeMux()
}

// Update the signature for the routes() so that it returnst a
// http.Handler instead of a *http.ServeMux
func (a *application) routes() http.Handler {
	// Use the http.NewServeMux() function to initialize a
	// new servemux, then register the home function as the
	// handler for the "/" URL pattern.
	mux.HandleFunc("/", a.home)
	mux.HandleFunc("/snippet", a.showSnippet)
	mux.HandleFunc("/snippet/create", a.createSnippet)

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
	mux.Handle("/static/", http.StripPrefix("/static", fs))

	// Pass the servemux as a 'next' parameter to the secureHeaders middleware.
	// Because secureHeaders is just a function, and the function returns a
	// http.Handler we don't need to do anything else.
	//Wrapping the existing chain with the logRequest middleware.
	return a.recoverPanic(a.logRequest(secureHeaders(mux)))
}
