package main

import (
	"flag"
	"log"
	"net/http"
	"os"
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

	// Create a new logger for writting information messages.
	infoLog := log.New(os.Stdout, "[INFO]\t", log.Ldate|log.Ltime)

	// Create a new logger for writting error messages. The log.Lshortfile flag to
	// include the relevant file name and line number
	errorLog := log.New(os.Stderr, "[ERROR]\t", log.Ldate|log.Ltime|log.Lshortfile)

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

	// Initialize a new http.Server struct. We set the Addr and Handler fields
	// that the server uses the same network address and routes as before, and
	// the ErrorLog field so that the server now uses the custom errorLog logger.
	svr := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  mux,
	}

	infoLog.Printf("Starting server on %s\n", *addr)

	// Call the ListenAndServe() method on our new http.Server struct.
	// If svr.ListenAndServe() returns an error we use the log.Fatal()
	// function to log the error message and exit.
	if err := svr.ListenAndServe(); err != nil {
		errorLog.Fatal(err)
	}
}
