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
		app.render(w, r, "login.page.tmpl", nil)
		return
	case "POST":
		err := r.ParseForm()
		if err != nil {
			app.serverError(w, err)
			return
		}
		form := forms.New(r.PostForm)
		_, err = app.models.Users.Authenticate(form.Get("username"), form.Get("password"))
		if err == data.ErrInvalidCredentials {
			form.Errors.Add("generic", "Email or Password is incorrect")
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
		app.render(w, r, "register.page.tmpl", nil)
		return
	case "POST":
		err := r.ParseForm()
		if err != nil {
			app.serverError(w, err)
			return
		}
		userData := forms.New(r.PostForm)
		userData.Required("username", "password", "confirm_password")
		userData.ConfirmPassword("password", "confirm_password")
		if !userData.Valid() {
			app.render(w, r, "register.page.tmpl", nil)
			return
		}
		_, err = app.models.Users.Insert(userData.Get("username"), userData.Get("password"))
		if err == data.ErrDuplicateUsername {
			userData.Errors.Add("username", "Username is already in use")
			app.render(w, r, "register.page.tmpl", nil)
			return
		} else if err != nil {
			app.serverError(w, err)
			return
		}

	}
	http.Redirect(w, r, "/", http.StatusSeeOther)

}
