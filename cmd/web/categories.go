package main

import (
	"net/http"
)

func (app *application) showCategory(w http.ResponseWriter, r *http.Request, category string) {

	// cookie, _ := r.Cookie("session")
	// fmt.Println(cookie)
	// fmt.Println("cookie")
	// fmt.Println(cookie.Expires)

	sortColumn := "created"
	time := "9999"
	// category := r.URL.Path[10:]
	// uq := r.URL.Query()
	// if uq.Get("col") != "" {
	// 	sortColumn = uq.Get("col")
	// }
	// if uq.Get("time") != "" {
	// 	time = uq.Get("time")
	// }

	categories, err := app.models.Categories.GetOne(category)

	if err != nil {
		app.serverError(w, err)
	}

	currentUser := app.contextGetUser(r)

	users, err := app.models.Users.GetAllUsers(currentUser.Name)

	if err != nil {
		app.serverError(w, err)
	}

	// if sortColumn == "top" {
	// 	posts, err := app.models.Posts.Latest(category, "votes", "DESC", time)
	// 	if err != nil {
	// 		app.serverError(w, err)
	// 	}
	// 	// app.render(w, r, "showcat.page.tmpl", &templateData{Posts: posts, Sort: "top", Categories: categories})
	// 	return
	// }
	posts, err := app.models.Posts.Latest(category, sortColumn, "DESC", time)
	if err != nil {
		app.serverError(w, err)
	}

	data := &templateData{Posts: posts, Categories: categories, Users: users}
	app.serveAsJSON(w, data)

	// j, err := json.Marshal(data)
	// if err != nil {
	// 	app.serverError(w, err)
	// }

	// io.Copy(w, bytes.NewReader(j))

	// app.render(w, r, "showcat.page.tmpl", &templateData{Posts: posts, Categories: categories})
}
