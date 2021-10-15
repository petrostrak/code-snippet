package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/petrostrak/code-snippet/pkg/models"
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

	s, err := a.snippets.Latest()
	if err != nil {
		a.serverError(w, err)
		return
	}

	// Use the new render helper.
	a.render(w, r, "home.page.tmpl", &templateData{
		Snippets: s,
	})
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

	s, err := a.snippets.Get(id)
	if err == models.ErrNoRecord {
		a.notFound(w)
		return
	} else if err != nil {
		a.serverError(w, err)
		return
	}

	// Create an instance of a templateData struct holding the snippet data.
	data := &templateData{Snippet: s}

	// Initialize a slice containing the paths to the show.page.tmpl file
	// plus the base layout and footer partial.
	files := []string{
		"./ui/html/show.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}

	// Parse the template files.
	ts, err := template.ParseFiles(files...)
	if err != nil {
		a.serverError(w, err)
		return
	}

	// And then execute them. Notice how we are passing in the templateData
	// struct as the final parameter.
	if err := ts.Execute(w, data); err != nil {
		a.serverError(w, err)
		return
	}

	// Use the fmt.Fprintf function to interpolate the id value with our
	// response and write it to the http.ResponseWriter.
	fmt.Fprintf(w, "%v", s)

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
	id, err := a.snippets.Insert(title, content, expires)
	if err != nil {
		a.serverError(w, err)
		return
	}

	// Redirect the user to the relevant page for the snippet.
	http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)
}
