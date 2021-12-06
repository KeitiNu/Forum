package main

import "net/http"

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// If a user enters a address that doesn't exist display 404 Error!
	if r.URL.Path != "/" {
		app.render(w, r, "400.page.tmpl", nil)
		return
	}
	app.render(w, r, "home.page.tmpl", nil)
}

func (app *application) login(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "login.page.tmpl", nil)
}
func (app *application) register(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "register.page.tmpl", nil)
}
