package main

import (
	"net/http"
	"time"

	"git.01.kood.tech/roosarula/forum/pkg/data"
	"git.01.kood.tech/roosarula/forum/pkg/forms"
	uuid "github.com/satori/go.uuid"
)

// All of the functions that run when a user enters an address are located here.
// In the next file "routes.go", you can observe which function is applied to which address.

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// If a user types in an address that doesn't exist, a 404 Error is displayed.
	if r.URL.Path != "/" {
		w.WriteHeader(404)
		app.render(w, r, "400.page.tmpl", nil)
		return
	}
	// If the user connects to the home address, the frontpage is displayed.
	app.render(w, r, "home.page.tmpl", nil)
}

func (app *application) login(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET": // When a person clicks the login link, the form appears.
		app.render(w, r, "login.page.tmpl", &templateData{Form: forms.New(nil)})
		return
	case "POST": // If a user submits a form on the login page, we check the data and then run the database queries.
		err := r.ParseForm()
		if err != nil {
			app.serverError(w, err)
			return
		}

		// We make a form object with user input and error storage.
		form := forms.New(r.PostForm)
		v := forms.NewValidator()
		form.Errors = v

		// User object
		user := &data.User{
			Name: form.Get("username"),
		}
		user.Password.Set(form.Get("password"))

		// Validate the user struct and return the error messages to the client if any of
		// the checks fail.
		if data.ValidateLogin(v, user); !v.Valid() {
			app.render(w, r, "login.page.tmpl", &templateData{Form: form})
			return
		}

		// Authenticate the user when the input is correct. If the credentials do not match, the user will receive a generic error message.
		// A generic error prevents from checking to see if an email address exists in our user database and start hacking.
		err = app.models.Users.Authenticate(user.Name, form.Get("password"))
		if err == data.ErrInvalidCredentials {
			form.Errors.AddError("generic", "Username or Password is incorrect")
			app.render(w, r, "login.page.tmpl", &templateData{Form: form})
			return
		} else if err != nil {
			app.serverError(w, err)
			return
		}

		// Get the token for the current user who is attempting to log in.
		a, err := r.Cookie("session")
		if err != nil {
			app.serverError(w, err)
		}

		// Add the current cookie (token) to the user's profile in database.
		err = app.models.Users.UpdateByToken(a.Value, user.Name)
		if err != nil {
			app.serverError(w, err)
			return
		}
	}
	// After login redirect the user to the homepage.
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
func (app *application) register(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		app.render(w, r, "register.page.tmpl", &templateData{Form: forms.New(nil)})
		return
	case "POST":
		err := r.ParseForm()
		if err != nil {
			app.serverError(w, err)
			return
		}

		form := forms.New(r.PostForm)
		v := forms.NewValidator()
		form.Errors = v

		user := &data.User{
			Name:  form.Get("username"),
			Email: form.Get("email"),
		}
		err = user.Password.Set(form.Get("password"))
		if err != nil {
			app.serverError(w, err)
			return
		}

		// Validate the user struct and return the error messages to the client if any of
		// the checks fail.
		if data.ValidateUser(v, user); !v.Valid() {
			app.render(w, r, "register.page.tmpl", &templateData{Form: form})
			return
		}

		// Get the token for the current user who is attempting to register.
		a, err := r.Cookie("session")
		if err != nil {
			app.serverError(w, err)
		}

		err = app.models.Users.Insert(user, a.Value)
		if err != nil {
			switch err {
			case data.ErrDuplicateUsername:
				form.Errors.AddError("username", "Username is already in use")
				app.render(w, r, "register.page.tmpl", &templateData{Form: form})
				return
			case data.ErrDuplicateEmail:
				form.Errors.AddError("email", "Email is already in use")
				app.render(w, r, "register.page.tmpl", &templateData{Form: form})
				return
			default:
				app.serverError(w, err)
				return
			}
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
}

func (app *application) logout(w http.ResponseWriter, r *http.Request) {
	c := &http.Cookie{Name: "session", Value: uuid.NewV4().String(), Expires: time.Now(), MaxAge: -1}
	http.SetCookie(w, c)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) profile(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "profile.page.tmpl", nil)
}
