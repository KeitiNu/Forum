package main

import (
	"fmt"
	"net/http"

	"git.01.kood.tech/roosarula/forum/pkg/data"
	"git.01.kood.tech/roosarula/forum/pkg/forms"
)

func (app *application) newCategory(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET": // When a person clicks the login link, the form appears.
		app.render(w, r, "createcat.page.tmpl", &templateData{Form: forms.New(nil)})
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

		category := &data.Category{
			Title:       form.Get("catname"),
			Description: form.Get("description"),
		}

		v.Check(category.Title != "", "title", "must be provided")
		v.Check(category.Description != "", "description", "must be provided")
		if !v.Valid() {
			app.render(w, r, "createcat.page.tmpl", &templateData{Form: form})
			return
		}

		app.models.Categories.Insert(category.Title, category.Description)
		if err != nil {
			fmt.Println("error happened", err)
		}
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) showCategory(w http.ResponseWriter, r *http.Request) {
	sortColumn := "created"
	category := r.URL.Path[10:]
	uq := r.URL.Query()
	if uq.Get("col") != "" {
		sortColumn = uq.Get("col")
	}
	posts, err := app.models.Posts.Latest(category, sortColumn, "DESC")
	if err != nil {
		app.serverError(w, err)
	}
	app.render(w, r, "showcat.page.tmpl", &templateData{Posts: posts})
}
