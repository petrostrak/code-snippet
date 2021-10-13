package main

import (
	"flag"
	"log"
	"net/http"
)

var (
	// http.ServeMux is also a handler, which instead of providing
	// a response itself passes the request on to a second handler.
	mux *http.ServeMux
)

func init() {
	mux = http.NewServeMux()
}

func StartApp() {

	// Define a new command-line flag with the name 'addr', a default value
	// and some sort help text explaining what the flag controls. The value
	// of the flag will be stored in the addr variable at runtime.
	addr := flag.String("addr", ":4000", "HTTP network address")

	// Importantly, we use the flag.Parse() to parse the command-line imput.
	flag.Parse()

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
	log.Printf("Starting server on %s\n", *addr)
	if err := http.ListenAndServe(*addr, mux); err != nil {
		log.Fatal(err)
	}
}
