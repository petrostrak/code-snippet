package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/justinas/nosurf"
	"github.com/petrostrak/code-snippet/pkg/models"
)

// Because we want secureHeaders middleware to act on every request that is
// received, we need it to be executed before a request hits our servemux.
// To do this, we need the secureHeaders middleware to wrap our servemux.
func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Frame-Options", "deny")

		next.ServeHTTP(w, r)
	})
}

func (a *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL)

		next.ServeHTTP(w, r)
	})
}

func (a *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a deferred function (which will always be run in the event
		// of a panic as Go unwinds the stack)
		defer func() {
			// Use the build-in recover function to check if there has been a
			// panic or not
			if err := recover(); err != nil {
				// Set a "Connection: close" header on the response.
				w.Header().Set("Connection", "close")
				// Call the a.serverError helper method to return a 500
				// Internal Server response.
				a.serverError(w, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (a *application) requireAuthenticatedUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// If the user is not authenticated, redirect them to the login page and
		// return from the middleware chain so that no subsequent handlers in the
		// chain are executed.
		if a.authenticatedUser(r) == nil {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}

		// Otherwise call the next handler in the chain.
		next.ServeHTTP(w, r)
	})
}

// Create a noSurf middleware function which uses a customized CSRF cookie with
// the Secure, Path and HttpOnly flags set.
func noSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	})

	return csrfHandler
}

func (a *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if a userID value exists in the session. If this isn't
		// present, then call the next handler in the chain as normal
		exists := a.session.Exists(r, "userID")
		if !exists {
			next.ServeHTTP(w, r)
			return
		}

		// Fetch the details of the current user from the database. If no
		// matching record is found, remove the (invalid) userID from their
		// session and call the next handler in the chain as normal.
		user, err := a.users.Get(a.session.GetInt(r, "userID"))
		if err == models.ErrNoRecord {
			a.session.Remove(r, "userID")
			next.ServeHTTP(w, r)
			return
		} else if err != nil {
			a.serverError(w, err)
			return
		}

		// Otherwise, we know that the request is coming from a valid
		// authenticated (logged in) user. We create a new copy of the
		// request with the user information added to the request context,
		// and call the next handle in the chain using this new copy of the
		// request.
		ctx := context.WithValue(r.Context(), contextKeyUser, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
