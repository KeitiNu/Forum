package main

import (
	"net/http"
)

func (app *application) showCategory(w http.ResponseWriter, r *http.Request, category string) {

	sortColumn := "created"
	time := "9999"

	categories, err := app.models.Categories.GetOne(category)

	if err != nil {
		app.serverError(w, err)
	}

	currentUser := app.contextGetUser(r)

	users, err := app.models.Users.GetAllUsers(currentUser.Name)

	if err != nil {
		app.serverError(w, err)
	}

	posts, err := app.models.Posts.Latest(category, sortColumn, "DESC", time)

	if err != nil {
		app.serverError(w, err)
	}

	data := &templateData{Posts: posts, Categories: categories, Users: users, AuthenticatedUser: currentUser}
	app.serveAsJSON(w, data)

}
