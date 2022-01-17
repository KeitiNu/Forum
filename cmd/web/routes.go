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

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/login", app.login)
	mux.HandleFunc("/logout", app.logout)
	mux.HandleFunc("/signup", app.register)
	mux.HandleFunc("/profile/", app.requireAuthenticatedUser(app.profile))
	mux.HandleFunc("/newcategory", app.requireAuthenticatedUser(app.newCategory))
	mux.HandleFunc("/submit", app.requireAuthenticatedUser(app.submitPost))
	mux.HandleFunc("/post/", app.showPost)
	mux.HandleFunc("/category/", app.showCategory)
	mux.HandleFunc("/edit/", app.editPost)
	mux.HandleFunc("/editimage/", app.editPostImage)
	mux.HandleFunc("/delete/", app.deletePost)
	mux.HandleFunc("/deletecomment/", app.deleteComment)
	mux.HandleFunc("/test", app.test)
	mux.HandleFunc("/testcomment", app.testcomment)

	return app.authenticate(app.session(mux))
}
