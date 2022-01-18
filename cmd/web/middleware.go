package main

import (
	"fmt"
	"net/http"
	"time"

	uuid "github.com/satori/go.uuid"
)

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a deferred function (which will always be run in the event
		// of a panic as Go unwinds the stack).
		defer func() {
			// Use the builtin recover function to check if there has been a
			// panic or not. If there has...
			if err := recover(); err != nil {
				// Set a "Connection: close" header on the response.
				w.Header().Set("Connection", "close")
				// Call the app.serverError helper method to return a 500
				// Internal Server response.
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// From the initial visit, the user will be given a session.
func (app *application) session(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := r.Cookie("session")
		if err != nil {
			c := &http.Cookie{Name: "session", Value: uuid.NewV4().String(), Expires: time.Now().Add(time.Hour), Path: "/"}
			http.SetCookie(w, c)
		}
		next.ServeHTTP(w, r)
	})
}

// The authentication token can be found in the cookies of your web browser.
func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a, err := r.Cookie("session")
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		// Let's locate the user linked with it if we get one.
		user, err := app.models.Users.GetByToken(a.Value)

		// Continue if there are no users associated with token.
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		// Otherwise, we can see that there is a match between the user in the database and cookie.
		// Add their information to the context of our request.
		r = app.contextSetUser(r, user)
		// Call the next handler in the chain.
		next.ServeHTTP(w, r)
	})
}

// This function is used to see if a user is logged in before they can publish or comment on our website.
func (app *application) requireAuthenticatedUser(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check that the authenticatedUser helper doesn't return nil.
		if app.contextGetUser(r) == nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		next.ServeHTTP(w, r)
	})
}
