package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"git.01.kood.tech/roosarula/forum/pkg/data"
	"git.01.kood.tech/roosarula/forum/pkg/forms"
)

type CommentForm struct {
	Comment string
	PostID  string
}

func (app *application) submitPost(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET": // When a person clicks the make a post link, the form appears.
		categories, err := app.models.Categories.Latest()
		if err != nil {
			fmt.Println(err)
		}
		// app.render(w, r, "submitpost.page.tmpl", &templateData{Form: forms.New(nil), Categories: categories})
		app.serveAsJSON(w, &templateData{Form: forms.New(nil), Categories: categories})
		return
	case "POST": // If a user submits a form on the login page, we check the data and then run the database queries.
		// var fileBytes []byte
		// err := r.ParseMultipartForm(20 << 20)
		// if err != nil {
		// 	app.serverError(w, err)
		// 	return
		// }
		// tempFilename := ""
		// file, _, err := r.FormFile("myFile")

		// // err.Error() asendus
		// if fmt.Sprintf("%s", err) != "http: no such file" {
		// 	if err != nil {
		// 		fmt.Println(err)
		// 		return
		// 	}
		// 	defer file.Close()

		// 	fileBytes, err = ioutil.ReadAll(file)
		// 	if err != nil {
		// 		fmt.Println(err)
		// 	}

		// 	if len(fileBytes) <= 20000000 {
		// 		tempFile, err := ioutil.TempFile("./ui/assets/thread-images", "upload-*.png")
		// 		if err != nil {
		// 			fmt.Println(err)
		// 			return
		// 		}

		// 		defer tempFile.Close()

		// 		_, err = tempFile.Write(fileBytes)
		// 		if err != nil {
		// 			app.serverError(w, err)
		// 		}
		// 		tempFilename = tempFile.Name()
		// 	}
		// }

		decoder := json.NewDecoder(r.Body)

		// We make a form object with user input and error storage.
		form := forms.New(r.PostForm)

		v := forms.NewValidator()
		form.Errors = v

		post := &data.Post{
			// Title:    form.Get("title"),
			// Content:  form.Get("content"),
			// Category: r.Form["category"],
		}

		err := decoder.Decode(&post)

		if err != nil {

			app.serverError(w, err)
			return
		}

		user := app.contextGetUser(r)
		post.User = user.Name
		categoryList, _ := app.models.Categories.Latest()
		v.Check(post.Title != "", "title", "Title must be provided")
		v.Check(post.Content != "", "content", "Description must be provided")
		v.Check(len(post.Category) != 0, "category", "At least 1 category must be provided")
		// v.Check(len(fileBytes) <= 20000000, "image", "Image can't be over 20 MB")
		// if !v.Valid() {
		// app.render(w, r, "submitpost.page.tmpl", &templateData{Form: form, Categories: categoryList})

		// return
		// }
		categories := post.Category
		id, err := app.models.Posts.Insert(post.Title, post.Content, post.User, post.ImageSrc, categories)
		if err != nil {
			fmt.Println("error happened", err)
		}

		users, err := app.models.Users.GetAllUsers()

		if err != nil {
			app.serverError(w, err)
		}

		stringId := strconv.Itoa(id)
		app.serveAsJSON(w, &templateData{Form: form, Categories: categoryList, Sort: stringId, Users: users})

		return
	}
}

func (app *application) showPost(w http.ResponseWriter, r *http.Request, idString string) {
	id, err := strconv.Atoi(idString)
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

	users, err := app.models.Users.GetAllUsers()

	if err != nil {
		app.serverError(w, err)
	}



	switch r.Method {
	case "POST":

		if user == nil {
			user = &data.User{}
		}

		data := &templateData{User: user, Post: post, Comments: comments, Users: users}

		j, err := json.Marshal(data)
		if err != nil {
			app.serverError(w, err)
		}

		io.Copy(w, bytes.NewReader(j))

		return

	case "GET":
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

func (app *application) comment(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		app.serverError(w, err)
		return
	}

	var c CommentForm

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&c)
	if err != nil {

		app.serverError(w, err)
		return
	}

	fmt.Println("C", c)

	user := app.contextGetUser(r)
	form := forms.New(r.PostForm)

	fmt.Println("FORM", form)

	v := forms.NewValidator()
	form.Errors = v

	id, err := strconv.Atoi(c.PostID)
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

	switch r.Method {
	case "GET":

		if user == nil {
			user = &data.User{}
		}

		data := &templateData{User: user, Post: post, Comments: comments}

		j, err := json.Marshal(data)
		if err != nil {
			app.serverError(w, err)
		}

		io.Copy(w, bytes.NewReader(j))

		return

	case "POST":
		fmt.Println(c.Comment)

		v.Check(c.Comment != "", "comment", "Cannot add empty comment")
		if !v.Valid() {
			app.serveAsJSON(w, &templateData{Form: form, User: user, Post: post, Comments: comments})
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
				Content: c.Comment,
			}
			user := app.contextGetUser(r)
			comment.User = user.Name

			_, err = app.models.Comments.Insert(comment)
			if err != nil {
				app.serverError(w, err)
			}
		}

		
	}

	
	users, err := app.models.Users.GetAllUsers()

	if err != nil {
		app.serverError(w, err)
	}
	app.serveAsJSON(w, &templateData{Form: form, User: user, Post: post, Comments: comments, Users: users})
	// http.Redirect(w, r, fmt.Sprintf("/post/%d", id), http.StatusFound)
}
