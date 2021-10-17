package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golangcollege/sessions"
	"github.com/petrostrak/code-snippet/pkg/models/mysql"
)

// Define an application struct to hold the application-wide dependencies for
// the web-app. Adding a snippet field to the struct will allow us to make the
// SnippetModel object available to our handlers
type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	session       *sessions.Session
	snippets      *mysql.SnippetModel
	templateCache map[string]*template.Template
}

func StartApp() {

	// Define a new command-line flag with the name 'addr', a default value
	// and some sort help text explaining what the flag controls. The value
	// of the flag will be stored in the addr variable at runtime.
	addr := flag.String("addr", ":4000", "HTTP network address")

	// Define a new command-line flag for the MySQL DSN string.
	dsn := flag.String("dsn", "ptrak:Password!@#$@/codesnippet?parseTime=true", "MySQL database")

	//Define a new command-line flag for the session secret (a random key which
	// will be used to encrypt and authenticate session cookies). It should be 32
	// bytes long.
	secret := flag.String("secret", "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "Secret")

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

	// Initialize a new template cache
	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	// Use the session.New() function to initialize a new session manager
	// passing in the secret key as the parameter. Then we configure it so
	// session always expires after 12 hours.
	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour
	session.Secure = true

	// Initialize a new instance of application containing the dependencies.
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		session:  session,
		// Initialize a mysql.SnippetModel instance and add it to the application
		// dependencies.
		snippets:      &mysql.SnippetModel{DB: db},
		templateCache: templateCache,
	}

	// Initialize a tls.Config struct to hold the non-default TLS settings we
	// will set the server to use.
	tlsConfig := &tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	// Initialize a new http.Server struct. We set the Addr and Handler fields
	// that the server uses the same network address and routes as before, and
	// the ErrorLog field so that the server now uses the custom errorLog logger.
	//
	// Set the server's TLSConfig field to use the tlsConfig variable we just
	// created.
	svr := &http.Server{
		Addr:      *addr,
		ErrorLog:  errorLog,
		Handler:   app.routes(),
		TLSConfig: tlsConfig,
		// Add Idle, Read and Write timeouts to the server
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("Starting server on %s\n", *addr)

	// Call the ListenAndServe() method on our new http.Server struct.
	// If svr.ListenAndServe() returns an error we use the log.Fatal()
	// function to log the error message and exit.
	//
	// Use the ListenAndServeTLS() to start the HTTP server. We
	// pass in the paths to the TLS certificate and corresponding private
	// key.
	if err := svr.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem"); err != nil {
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
