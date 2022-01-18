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

		v.Check(category.Title != "", "title", "Name must be provided")
		v.Check(category.Description != "", "description", "Description must be provided")
		if !v.Valid() {
			app.render(w, r, "createcat.page.tmpl", &templateData{Form: form})
			return
		}

		err = app.models.Categories.Insert(category.Title, category.Description)
		if err != nil {
			if err.Error() == "UNIQUE constraint failed: categories.title" {
				v.Check(1 > 2, "title", "Category already exists")
				if !v.Valid() {
					app.render(w, r, "createcat.page.tmpl", &templateData{Form: form})
					return
				}
			}
			fmt.Println("error happened", err)
		}
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) showCategory(w http.ResponseWriter, r *http.Request) {
	sortColumn := "created"
	time := "9999"
	category := r.URL.Path[10:]
	uq := r.URL.Query()
	if uq.Get("col") != "" {
		sortColumn = uq.Get("col")
	}
	if uq.Get("time") != "" {
		time = uq.Get("time")
	}

	categories, err := app.models.Categories.GetOne(category)

	if err != nil {
		app.serverError(w, err)
	}

	if sortColumn == "top" {
		posts, err := app.models.Posts.Latest(category, "votes", "DESC", time)
		if err != nil {
			app.serverError(w, err)
		}
		app.render(w, r, "showcat.page.tmpl", &templateData{Posts: posts, Sort: "top", Categories: categories})
		return
	}
	posts, err := app.models.Posts.Latest(category, sortColumn, "DESC", time)
	if err != nil {
		app.serverError(w, err)
	}
	app.render(w, r, "showcat.page.tmpl", &templateData{Posts: posts, Categories: categories})
}
