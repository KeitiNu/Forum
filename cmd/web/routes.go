package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	// We create a fileserver and strip the prefix for asset files. They are in /ui/assets/ folder but we present them
	// as they would be in a folder called /static/. This gives us some security and that means that our calls for the files
	// in code do not need to match the computer filesystem.

	fs := http.FileServer(http.Dir("./ui/assets"))
	mux.Handle("/static/", http.StripPrefix("/static", fs))

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/login", app.login)
	mux.HandleFunc("/logout", app.logout)
	mux.HandleFunc("/register", app.register)
	return app.authenticate(app.session(mux))
}
