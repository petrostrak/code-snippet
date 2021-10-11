package main

import (
	"log"
	"net/http"
)

// Define a home hundler function which writes a byte of
// slice containing "Hello from Code Snippet!" as the
// response body.
func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from Code Snippet!"))
}

func main() {

	// Use the http.NewServeMux() function to initialize a
	// new servemux, then register the home function as the
	// handler for the "/" URL pattern.
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)

	// Use the http.ListenAndServe() function to start a new web server. We pass
	// two parameters: the TCP network address to listen on (in this case ":4000)
	// and the servemux we just created. If http.ListenAndServe() returns an error
	// we use the log.Fatal() function to log the error message and exit.
	log.Println("Starting server on :4000")
	if err := http.ListenAndServe(":4000", mux); err != nil {
		log.Fatal(err)
	}
}
