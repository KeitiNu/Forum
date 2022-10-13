package main

import (
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	// We create a fileserver and strip the prefix for asset files. They are in /ui/assets/ folder but we present them
	// as they would be in a folder called /static/. This gives us some security and that means that our calls for the files
	// in code do not need to match the computer filesystem.

	fs := http.FileServer(http.Dir("./ui/assets"))
	mux.Handle("/static/", http.StripPrefix("/static", fs))
	mux.Handle("/category/static/", http.StripPrefix("/category/static", fs))
	mux.Handle("/post/static/", http.StripPrefix("/post/static", fs))

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/socket", app.socket)
	mux.HandleFunc("/data/", app.data)

	// mux.HandleFunc("/post/", app.showPost)
	// mux.HandleFunc("/category/", app.showCategory)

	// mux.HandleFunc("/login", app.login)

	mux.HandleFunc("/logout", app.logout)
	mux.HandleFunc("/message", app.message)
	mux.HandleFunc("/typing", app.typing)

	// mux.HandleFunc("/signup", app.register)

	// mux.HandleFunc("/profile/", app.requireAuthenticatedUser(app.profile))
	// mux.HandleFunc("/newcategory", app.requireAuthenticatedUser(app.newCategory))
	// mux.HandleFunc("/submit", app.requireAuthenticatedUser(app.submitPost))

	return app.recoverPanic(app.authenticate(app.session(mux)))
}
