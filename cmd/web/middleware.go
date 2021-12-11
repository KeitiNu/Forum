package main

import (
	"net/http"
	"time"

	uuid "github.com/satori/go.uuid"
)

/* User will get a session from first visit.
If cookie doesnt exits assign one, else continue */
/* If user logs in add this token to this user
If there is a user with this token then authenticate */
func (app *application) session(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := r.Cookie("session")
		if err != nil {
			c := &http.Cookie{Name: "session", Value: uuid.NewV4().String(), Expires: time.Now().Add(time.Hour)}
			http.SetCookie(w, c)
		}
		next.ServeHTTP(w, r)
	})
}

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a, err := r.Cookie("session")
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		user, err := app.models.Users.GetByToken(a.Value)

		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		// Call the contextSetUser() helper to add the user information to the request
		// context.
		r = app.contextSetUser(r, user)
		// Call the next handler in the chain.
		next.ServeHTTP(w, r)
	})
}

func (app *application) requireAuthenticatedUser(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check that the authenticatedUser helper doesn't return nil.
		if app.contextGetUser(r) == nil {
			http.Redirect(w, r, "/login", 302)
			return
		}

		next.ServeHTTP(w, r)
	})
}
