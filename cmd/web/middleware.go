package main

import (
	"fmt"
	"net/http"

	"github.com/justinas/nosurf"
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
		if a.authenticatedUser(r) == 0 {
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
