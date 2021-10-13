package main

import (
	"log"
	"net/http"
)

var (
	mux *http.ServeMux
)

func init() {
	mux = http.NewServeMux()
}

func StartApp() {
	// Use the http.NewServeMux() function to initialize a
	// new servemux, then register the home function as the
	// handler for the "/" URL pattern.
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)

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

	// Use the http.ListenAndServe() function to start a new web server. We pass
	// two parameters: the TCP network address to listen on (in this case ":4000)
	// and the servemux we just created. If http.ListenAndServe() returns an error
	// we use the log.Fatal() function to log the error message and exit.
	log.Println("Starting server on :4000")
	if err := http.ListenAndServe(":4000", mux); err != nil {
		log.Fatal(err)
	}
}
