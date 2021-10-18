package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/petrostrak/code-snippet/pkg/forms"
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

	// Create a new forms.Form struct containing the POSTed date from the
	// form, then use the validation methods to check the content.
	form := forms.New(r.PostForm)
	form.Required("title", "content", "expires")
	form.MaxLength("title", 100)
	form.PermittedValues("expires", "365", "7", "1")

	// If the form isn't valid, redisplay the template passing in the
	// form.Form object as the data.
	if !form.Valid() {
		a.render(w, r, "create.page.tmpl", &templateData{Form: form})
		return
	}

	// Pass the data to the SnippetModel.Insert() receiving the ID of the new record back.
	id, err := a.snippets.Insert(form.Get("title"), form.Get("content"), form.Get("expires"))
	if err != nil {
		a.serverError(w, err)
		return
	}

	// Use the Put() method to add a string value ("Your snippet was saved
	// successfully!") and the corresponding key ("flash") to the session
	// data. Note that if there's no existing session for the current user
	// (or their session has expired) then a new, empty session for them
	// will automatically be created by the session middleware.
	a.session.Put(r, "flash", "Snippet successfully created!")

	// Redirect the user to the relevant page for the snippet.
	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}

func (a *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
	a.render(w, r, "create.page.tmpl", &templateData{
		// Pass a new empty forms.Form object to the template.
		Form: forms.New(nil),
	})
}

func (a *application) signupUserForm(w http.ResponseWriter, r *http.Request) {
	a.render(w, r, "signup.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func (a *application) signupUser(w http.ResponseWriter, r *http.Request) {
	// Parse the form data.
	if err := r.ParseForm(); err != nil {
		a.clientError(w, http.StatusBadRequest)
		return
	}

	// Validate the form contents using the form helper
	form := forms.New(r.PostForm)
	form.Required("name", "email", "password")
	form.MatchesPattern("email", forms.EmailRX)
	form.MinLength("password", 10)

	// If there are any errors, redisplay the signup form.
	if !form.Valid() {
		a.render(w, r, "signup.page.tmpl", &templateData{
			Form: form,
		})
		return
	}

	// Otherwise send a placeholder response
	fmt.Fprintln(w, "Create a new user")
}

func (a *application) loginUserForm(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Display the user login form")
}

func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Authenticate and login the user...")
}

func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Logout the user...")
}
