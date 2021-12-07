package main

import (
	"net/http"

	"git.01.kood.tech/roosarula/forum/pkg/data"
	"git.01.kood.tech/roosarula/forum/pkg/forms"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// If a user enters a address that doesn't exist display 404 Error!
	if r.URL.Path != "/" {
		app.render(w, r, "400.page.tmpl", nil)
		return
	}
	app.render(w, r, "home.page.tmpl", nil)
}

func (app *application) login(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		app.render(w, r, "login.page.tmpl", &templateData{Form: forms.New(nil)})
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
			Name: form.Get("username"),
		}
		user.Password.Set(form.Get("password"))

		// Validate the user struct and return the error messages to the client if any of
		// the checks fail.
		if data.ValidateUser(v, user); !v.Valid() {
			app.render(w, r, "login.page.tmpl", &templateData{Form: form})
			return
		}
		if err == data.ErrInvalidCredentials {
			form.Errors.AddError("generic", "Email or Password is incorrect")
			app.render(w, r, "login.page.tmpl", &templateData{Form: form})
			return
		} else if err != nil {
			app.serverError(w, err)
			return
		}
	}
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

		user := &data.User{
			Name:  form.Get("username"),
			Email: form.Get("email"),
		}
		err = user.Password.Set(form.Get("password"))
		if err != nil {
			app.serverError(w, err)
			return
		}

		v := forms.NewValidator()
		form.Errors = v
		// Validate the user struct and return the error messages to the client if any of
		// the checks fail.
		if data.ValidateUser(v, user); !v.Valid() {
			app.render(w, r, "register.page.tmpl", &templateData{Form: form})
			return
		}

		err = app.models.Users.Insert(user)
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
