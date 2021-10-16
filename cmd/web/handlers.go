package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/petrostrak/code-snippet/pkg/models"
)

// Define a home hundler function which writes a byte of
// slice containing "Hello from Code Snippet!" as the
// response body.
func (a *application) home(w http.ResponseWriter, r *http.Request) {
	// Because Pat matches the "/" path exactly, we can now remove the manual condition
	// of r.URL.Path != "/" from this handler.

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
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
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

	// Use the render helper
	a.render(w, r, "show.page.tmpl", &templateData{
		Snippet: s,
	})

}

// Add a createSnippet handler function .
// curl -i -X POST http://localhost:8080/snippet/create
func (a *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	// The check of r.Method != "POST" is now superfluous and can be removed.

	if err := r.ParseForm(); err != nil {
		a.clientError(w, http.StatusBadRequest)
		return
	}

	// Use the r.PostForm.Get() method to retrieve the relevant data fields
	// from the r.PostForm map.
	//
	// The r.PostForm map is populated only for POST, PATCH and PUT
	// requests, and contains the form data from the request body.
	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")
	expires := r.PostForm.Get("expires")

	// Initialize a map to hold any validation errors.
	//
	// More code patterns and validation visit https://www.alexedwards.net/blog/validation-snippets-for-go
	errors := make(map[string]string)

	// Check that the title field is not blank and not more that 100 characters
	// long. If it fails either of those checks, add a message to the errors
	// map using the field name as the key.
	if strings.TrimSpace(title) == "" {
		errors["title"] = "This field cannot be blank"
	} else if utf8.RuneCountInString(title) > 100 {
		errors["title"] = "This field is too long (maximum is 100 characters)"
	}

	// Check that the content field isn't blank.
	if strings.TrimSpace(content) == "" {
		errors["content"] = "This field cannot be blank"
	}

	// Check the expires field isn't blank and matches one of the permitted
	// values ("1", "7" or "365")
	if strings.TrimSpace(expires) == "" {
		errors["expires"] = "This field cannot be blank"
	} else if expires != "365" && expires != "7" && expires != "1" {
		errors["expires"] = "This field is invalid"
	}

	// If there are any errors, dump them in a plain text HTTP response and return
	// from the handler.
	if len(errors) > 0 {
		a.render(w, r, "create.page.tmpl", &templateData{
			FormErrors: errors,
			FormData:   r.PostForm,
		})
		return
	}

	// Pass the data to the SnippetModel.Insert() receiving the ID of the new record back.
	id, err := a.snippets.Insert(title, content, expires)
	if err != nil {
		a.serverError(w, err)
		return
	}

	// Redirect the user to the relevant page for the snippet.
	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}

func (a *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
	a.render(w, r, "create.page.tmpl", nil)
}
