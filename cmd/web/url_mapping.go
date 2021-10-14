package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/petrostrak/code-snippet/pkg/models/mysql"
)

// Define an application struct to hold the application-wide dependencies for
// the web-app. Adding a snippet field to the struct will allow us to make the
// SnippetModel object available to our handlers
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	snippet  *mysql.SnippetModel
}

func StartApp() {

	// Define a new command-line flag with the name 'addr', a default value
	// and some sort help text explaining what the flag controls. The value
	// of the flag will be stored in the addr variable at runtime.
	addr := flag.String("addr", ":4000", "HTTP network address")

	// Define a new command-line flag for the MySQL DSN string.
	dsn := flag.String("dsn", "ptrak:Password!@#$@/codesnippet?parseTime=true", "MySQL database")
	// Importantly, we use the flag.Parse() to parse the command-line imput.
	flag.Parse()

	// Create a new logger for writting information messages.
	infoLog := log.New(os.Stdout, "[INFO]\t", log.Ldate|log.Ltime)

	// Create a new logger for writting error messages. The log.Lshortfile flag to
	// include the relevant file name and line number
	errorLog := log.New(os.Stderr, "[ERROR]\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	// We also defer a call to db.Close() so that the connection pool is closed
	// before the main() returns.
	defer db.Close()

	// Initialize a new instance of application containing the dependencies.
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		// Initialize a mysql.SnippetModel instance and add it to the application
		// dependencies.
		snippet: &mysql.SnippetModel{DB: db},
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

// The openDB() function wraps sql.Open() and returns an sql.DB connection pool
// for a given DSN
func openDB(dsn string) (*sql.DB, error) {

	// The sql.Open() function doesn’t actually create any connections, all
	// it does is initialize the pool for future use. Actual connections to the
	// database are established lazily.
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
