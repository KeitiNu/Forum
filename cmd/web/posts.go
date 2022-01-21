package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

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
		var fileBytes []byte
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

			fileBytes, err = ioutil.ReadAll(file)
			if err != nil {
				fmt.Println(err)
			}
			
			if len(fileBytes) <= 20000000 {
				tempFile, err := ioutil.TempFile("./ui/assets/thread-images", "upload-*.png")
				if err != nil {
					fmt.Println(err)
					return
				}

				defer tempFile.Close()

				_, err = tempFile.Write(fileBytes)
				if err != nil {
					app.serverError(w, err)
				}
				tempFilename = tempFile.Name()
			}
		}

		// We make a form object with user input and error storage.
		form := forms.New(r.PostForm)
		v := forms.NewValidator()
		form.Errors = v

		post := &data.Post{
			Title:    form.Get("title"),
			Content:  form.Get("content"),
			Category: r.Form["category"],
			ImageSrc: tempFilename,
		}

		user := app.contextGetUser(r)
		post.User = user.Name
		categoryList, _ := app.models.Categories.Latest()
		v.Check(post.Title != "", "title", "Title must be provided")
		v.Check(post.Content != "", "content", "Description must be provided")
		v.Check(len(post.Category) != 0, "category", "At least 1 category must be provided")
		v.Check(len(fileBytes) <= 20000000, "image", "Image can't be over 20 Megabytes")
		if !v.Valid() {
			app.render(w, r, "submitpost.page.tmpl", &templateData{Form: form, Categories: categoryList})
			return
		}
		categories := r.Form["category"]
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
		return
	}
	post, err := app.models.Posts.Get(id)
	if err != nil {
		app.notFound(w)
		return
	}
	if post.Title == "" {
		app.notFound(w)
		return
	}
	comments, err := app.models.Comments.Latest(id)
	if err != nil {
		fmt.Println(err)
	}
	user := app.contextGetUser(r)
	switch r.Method {
	case "GET":

		if user == nil {
			user = &data.User{}
		}
		app.render(w, r, "showpost.page.tmpl", &templateData{User: user, Post: post, Comments: comments})
		return
	case "POST":
		err := r.ParseForm()
		if err != nil {
			app.serverError(w, err)
			return
		}
		user := app.contextGetUser(r)
		form := forms.New(r.PostForm)
		v := forms.NewValidator()
		form.Errors = v
		comment := form.Get("comment")
		v.Check(comment != "", "comment", "Cannot add empty comment")
		if !v.Valid() {
			app.render(w, r, "showpost.page.tmpl", &templateData{Form: form, User: user, Post: post, Comments: comments})
			return
		}
		if a := form.Get("commentUpdate"); a != "" {
			cid, _ := strconv.Atoi(form.Get("commentUpdateID"))
			if user.Name != form.Get("commentUpdateUser") {
				http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
				return
			}
			err = app.models.Comments.Update(a, cid)
			if err != nil {
				app.serverError(w, err)
			}
		} else {
			comment := &data.Comment{
				PostID:  id,
				Content: form.Get("comment"),
			}
			user := app.contextGetUser(r)
			comment.User = user.Name

			_, err = app.models.Comments.Insert(comment)
			if err != nil {
				app.serverError(w, err)
			}
		}
	}
	http.Redirect(w, r, fmt.Sprintf("/post/%d", id), http.StatusFound)
}

func (app *application) editPost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Path[6:])
	if err != nil {
		app.serverError(w, err)
		return
	}
	user := app.contextGetUser(r)
	switch r.Method {
	case "GET":
		post, err := app.models.Posts.Get(id)
		if user.Name != post.User {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}
		if err != nil {
			fmt.Println(err)
		}
		app.render(w, r, "editpost.page.tmpl", &templateData{Post: post})
	case "POST":
		err := r.ParseForm()
		if err != nil {
			app.serverError(w, err)
			return
		}
		form := forms.New(r.PostForm)
		v := forms.NewValidator()
		form.Errors = v

		post := &data.Post{
			Title:   form.Get("title"),
			Content: form.Get("content"),
			ID:      id,
		}
		err = app.models.Posts.Update(post.Title, post.Content, post.ID)
		if err != nil {
			fmt.Println(err)
		}
		http.Redirect(w, r, fmt.Sprintf("/post/%d", id), http.StatusFound)
		return
	}

}

func (app *application) editPostImage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("running editPostImage")
	id, err := strconv.Atoi(strings.Split(r.URL.Path[11:], "?")[0])
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = r.ParseForm()
	if err != nil {
		app.serverError(w, err)
		return
	}
	form := forms.New(r.PostForm)
	v := forms.NewValidator()
	form.Errors = v

	post := &data.Post{
		ImageSrc: form.Get("image"),
		ID:       id,
	}
	err = app.models.Posts.UpdateImage(post.ImageSrc, post.ID)
	if err != nil {
		fmt.Println(err)
	}
	http.Redirect(w, r, fmt.Sprintf("/edit/%d", id), http.StatusFound)
}

func (app *application) deletePost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Path[8:])
	if err != nil {
		app.serverError(w, err)
		return
	}
	user := app.contextGetUser(r)
	post, _ := app.models.Posts.Get(id)
	if user.Name != post.User {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}
	err = app.models.Posts.Delete(id)
	if err != nil {
		app.serverError(w, err)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func (app *application) deleteComment(w http.ResponseWriter, r *http.Request) {
	back := r.Header.Get("referer")
	id, err := strconv.Atoi(r.URL.Path[15:])
	if err != nil {
		app.serverError(w, err)
		return
	}
	user := app.contextGetUser(r)
	comment, _ := app.models.Comments.Get(id)
	if user.Name != comment.User {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}
	err = app.models.Comments.Delete(id)
	if err != nil {
		app.serverError(w, err)
		return
	}
	http.Redirect(w, r, back, http.StatusFound)

}

func (app *application) test(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(1 << 20)
	if err != nil {
		app.serverError(w, err)
	}

	user := app.contextGetUser(r)
	id := r.Form.Get("postID")
	vote := r.Form.Get("type")
	err = app.models.Posts.AddVote(id, vote, user.Name)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) testcomment(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(1 << 20)
	if err != nil {
		app.serverError(w, err)
	}
	user := app.contextGetUser(r)
	id := r.Form.Get("postID")
	vote := r.Form.Get("type")
	err = app.models.Comments.AddVote(id, vote, user.Name)
	if err != nil {
		app.serverError(w, err)
	}
}
