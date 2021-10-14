package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

// Define an application struct to hold the application-wide dependencies for
// the web-app.
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
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

	// Initialize a new instance of application containing the dependencies.
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	// Initialize a new http.Server struct. We set the Addr and Handler fields
	// that the server uses the same network address and routes as before, and
	// the ErrorLog field so that the server now uses the custom errorLog logger.
	svr := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on %s\n", *addr)

	// Call the ListenAndServe() method on our new http.Server struct.
	// If svr.ListenAndServe() returns an error we use the log.Fatal()
	// function to log the error message and exit.
	if err := svr.ListenAndServe(); err != nil {
		errorLog.Fatal(err)
	}
}
