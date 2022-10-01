package main

import (
	"context"
	"fmt"
	"net/http"

	"git.01.kood.tech/roosarula/forum/pkg/data"
)

// Define a custom contextKey type, with the underlying type string.
type contextKey string

// Convert the string "user" to a contextKey type and assign it to the userContextKey
// constant. We'll use this constant as the key for getting and setting user information
// in the request context.
const userContextKey = contextKey("users")

// The contextSetUser() method returns a new copy of the request with the provided
// User struct added to the context. Note that we use our userContextKey constant as the
// key.
func (app *application) contextSetUser(r *http.Request, user *data.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

// The contextSetUser() retrieves the User struct from the request context.
func (app *application) contextGetUser(r *http.Request) *data.User {

	var ctx = r.Context()
	fmt.Println(ctx)

	user, ok := r.Context().Value(userContextKey).(*data.User)

	if !ok {
		return nil
	}

	return user
}
