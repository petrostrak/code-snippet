package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

// Define a home hundler function which writes a byte of
// slice containing "Hello from Code Snippet!" as the
// response body.
func (a *application) home(w http.ResponseWriter, r *http.Request) {

	// Check if the current request URL path exaclty matches "/".
	// If it doesn't, the http.NotFound() function triggers to send
	// a 404 response to the client. Then we return to avoid executing
	// any following code.
	if r.URL.Path != "/" {
		a.notFound(w)
		return
	}

	// Initialize a slice containing the paths to the two files. Not that
	// the home.page.tmpl file must be the *first* file in the slice.
	files := []string{
		"./ui/html/home.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}

	// Use the template.ParseFiles() function to read the template file
	// into a template set. Notice that we can pass the slice of files as
	// a variadic parameter.
	ts, err := template.ParseFiles(files...)
	if err != nil {
		a.serverError(w, err)
		return
	}

	// We then use the Execute() method on the template set to write
	// the template content as the respose body. The last parameter to
	// Execute represents dynamic data that we want to pass in.
	if err := ts.Execute(w, nil); err != nil {
		a.serverError(w, err)
	}

}

// Add a showSnippet handler function.
func (a *application) showSnippet(w http.ResponseWriter, r *http.Request) {

	// Extract the value of the id parameter from the query string
	// and try to convert it to an integer using the strconv.Atoi()
	// function. If it cannot be converted to an integer of the value
	// is less that 1, we return a 404 not found response.
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		a.notFound(w)
		return
	}

	// Use the fmt.Fprintf function to interpolate the id value with our
	// response and write it to the http.ResponseWriter.
	fmt.Fprintf(w, "Display a specific snippet with id %d\n", id)
	w.Write([]byte("Display a specific snippet!"))
}

// Add a createSnippet handler function .
// curl -i -X POST http://localhost:8080/snippet/create
func (a *application) createSnippet(w http.ResponseWriter, r *http.Request) {

	// Use r.Method to check whether the request is using POST or not.
	// If it's not, use the w.WriteHeader() method to send a 405 status
	// code, the w.Write() method to write a response body and then
	// return from the function.
	if r.Method != http.MethodPost {

		// Use the Header().Set() method to add an 'Allow: Post' header to
		// the response header map. The first parameter is the header name
		// and the second parameter is the header value.
		w.Header().Set("Allow", http.MethodPost)

		// The clientError helper sends a specific status code and corresponding
		// description to the user.
		a.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	// Some dummy data
	title := "0 snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi"
	expires := "7"

	// Pass the data to the SnippetModel.Insert() receiving the ID of the new record back.
	id, err := a.snippet.Insert(title, content, expires)
	if err != nil {
		a.serverError(w, err)
		return
	}

	// Redirect the user to the relevant page for the snippet.
	http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)
}
