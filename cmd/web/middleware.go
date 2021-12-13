package main

import (
	"expvar"
	"net/http"
	"time"

	uuid "github.com/satori/go.uuid"
)

// From the initial visit, the user will be given a session.
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
			http.Redirect(w, r, "/login", 302)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (app *application) metrics(next http.Handler) http.Handler {
	// Initialize the new expvar variables when the middleware chain is first built.
	totalRequestsReceived := expvar.NewInt("total_requests_received")
	totalResponsesSent := expvar.NewInt("total_responses_sent")
	totalProcessingTimeMicroseconds := expvar.NewInt("total_processing_time_μs")

	// The following code will be run for every request...
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Record the time that we started to process the request.
		start := time.Now()

		// Use the Add() method to increment the number of requests received by 1.
		totalRequestsReceived.Add(1)

		// Call the next handler in the chain.
		next.ServeHTTP(w, r)

		// On the way back up the middleware chain, increment the number of responses
		// sent by 1.
		totalResponsesSent.Add(1)

		// Calculate the number of microseconds since we began to process the request,
		// then increment the total processing time by this amount.
		duration := time.Now().Sub(start).Microseconds()
		totalProcessingTimeMicroseconds.Add(duration)
	})
}
