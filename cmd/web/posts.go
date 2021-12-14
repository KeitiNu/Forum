package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"git.01.kood.tech/roosarula/forum/pkg/data"
	"git.01.kood.tech/roosarula/forum/pkg/forms"
)

func (app *application) submitPost(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET": // When a person clicks the make a post link, the form appears.
		categories, err := app.models.Categories.Latest()
		if err != nil {
			fmt.Println(err)
		}
		app.render(w, r, "submitpost.page.tmpl", &templateData{Form: forms.New(nil), Categories: categories})
		return
	case "POST": // If a user submits a form on the login page, we check the data and then run the database queries.
		err := r.ParseMultipartForm(20 << 20)
		if err != nil {
			app.serverError(w, err)
			return
		}
		tempFilename := ""
		file, _, err := r.FormFile("myFile")
		// err.Error() asendus
		if fmt.Sprintf("%s", err) != "http: no such file" {
			if err != nil {
				fmt.Println(err)
				return
			}
			defer file.Close()

			tempFile, err := ioutil.TempFile("./ui/assets/thread-images", "upload-*.png")
			if err != nil {
				fmt.Println(err)
				return
			}

			defer tempFile.Close()

			fileBytes, err := ioutil.ReadAll(file)
			if err != nil {
				fmt.Println(err)
			}

			tempFile.Write(fileBytes)
			tempFilename = tempFile.Name()
		}

		// We make a form object with user input and error storage.
		form := forms.New(r.PostForm)
		v := forms.NewValidator()
		form.Errors = v

		post := &data.Post{
			Title:    form.Get("title"),
			Content:  form.Get("content"),
			ImageSrc: tempFilename,
		}

		user := app.contextGetUser(r)
		post.User = user.Name
		v.Check(post.Title != "", "title", "must be provided")
		v.Check(post.Content != "", "content", "must be provided")
		if !v.Valid() {
			app.render(w, r, "submitpost.page.tmpl", &templateData{Form: form})
			return
		}
		categories := r.Form["category"]
		fmt.Println(post.ImageSrc)
		id, err := app.models.Posts.Insert(post.Title, post.Content, post.User, post.ImageSrc, categories)
		if err != nil {
			fmt.Println("error happened", err)
		}
		http.Redirect(w, r, fmt.Sprintf("/post/%d", id), http.StatusSeeOther)
		return
	}
}

func (app *application) showPost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Path[6:])
	if err != nil {
		app.serverError(w, err)
	}
	post, err := app.models.Posts.Get(id)
	comments, err := app.models.Comments.Latest(id)
	if err != nil {
		fmt.Println(err)
	}
	switch r.Method {
	case "GET":
		app.render(w, r, "showpost.page.tmpl", &templateData{Post: post, Comments: comments})
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

		comment := &data.Comment{
			PostID:  id,
			Content: form.Get("comment"),
		}
		user := app.contextGetUser(r)
		comment.User = user.Name

		app.models.Comments.Insert(comment)
		http.Redirect(w, r, fmt.Sprintf("/post/%d", id), 302)
		return
	}

}
