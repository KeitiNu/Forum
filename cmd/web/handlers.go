package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"git.01.kood.tech/roosarula/forum/pkg/data"
	"git.01.kood.tech/roosarula/forum/pkg/forms"
	uuid "github.com/satori/go.uuid"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		categories, err := app.models.Categories.Latest()
		if err != nil {
			app.serverError(w, err)
		}
		// currentUser := app.contextGetUser(r)

		// if err != nil {
		// 	app.serverError(w, err)
		// }
		// users, err := app.models.Users.GetAllUsers(currentUser.Name)

		app.render(w, r, "index.html", &templateData{Categories: categories})

	case "POST":
		app.serverError(w, errors.New("POST METHOD NOT ALLOWED"))
	}
}

func (app *application) message(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		err := r.ParseForm()
		if err != nil {
			app.serverError(w, err)
			return
		}

		var c forms.ChatBoxForm

		decoder := json.NewDecoder(r.Body)
		err = decoder.Decode(&c)
		if err != nil {

			app.serverError(w, err)
			return
		}

		messages, err := app.models.Messages.GetMessages(c.User, c.Recipient, c.Offset)

		if err != nil {

			app.serverError(w, err)
			return
		}

		j, err := json.Marshal(messages)
		if err != nil {
			app.serverError(w, err)
		}

		io.Copy(w, bytes.NewReader(j))

	case "GET":
		app.serverError(w, errors.New("POST METHOD NOT ALLOWED"))
	}
}

func (app *application) data(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		path := r.URL.Path[6:]
		paths := strings.Split(path, "/")

		fmt.Println("PATH: ", path)

		switch paths[0] {
		case "category":
			app.showCategory(w, r, paths[1])
		case "post":
			app.showPost(w, r, paths[1])
		case "login":
			app.login(w, r)
		case "submit":
			app.submitPost(w, r)
		case "signup":
			app.register(w, r)
		case "profile":
			app.profile(w, r)
		case "comment":
			app.comment(w, r)
		case "chat":
			app.chat(w, r)
		default:
			categories, _ := app.models.Categories.Latest()
			currentUser := app.contextGetUser(r)
			users, _ := app.models.Users.GetAllUsers(currentUser.Name)

			app.serveAsJSON(w, &templateData{Categories: categories, AuthenticatedUser: currentUser, Users: users})

			// app.serveAsJSON(w, &templateData{Categories: categories})

		}
	case "GET":
		app.serverError(w, errors.New("GET METHOD NOT ALLOWED"))
	}

}

func (app *application) chat(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.serverError(w, err)
		return
	}

	var c forms.ChatForm

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&c)
	if err != nil {

		app.serverError(w, err)
		return
	}

	var msg = &data.Message{
		Recipient: c.RecipientId,
		Sender:    c.UserId,
		Content:   c.Message,
	}

	err = app.models.Messages.Insert(msg)
	if err != nil {
		app.serverError(w, err)
		return
	}

	sendChatNotification(msg.Sender, msg.Recipient, msg.Content)
}

func (app *application) login(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET": // When a person clicks the login link, the form appears.
		app.home(w, r)
	case "POST": // If a user submits a form on the login page, we check the data and then run the database queries.
		var p forms.LoginForm
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&p)
		if err != nil {

			app.serverError(w, err)
			return
		}

		// err := r.ParseForm()

		// if err != nil {

		// 	app.serverError(w, err)
		// 	return
		// }

		// We make a form object with user input and error storage.
		form := forms.New(r.PostForm)

		v := forms.NewValidator()
		form.Errors = v

		// User object
		user := &data.User{
			Name: p.Username,
		}
		err = user.Password.Set(p.Password)

		if err != nil {
			app.serverError(w, err)
		}

		// Validate the user struct and return the error messages to the client if any of
		// the checks fail.
		if data.ValidateLogin(v, user); !v.Valid() {
			// app.render(w, r, "index.html", &templateData{Form: form})
			app.serveAsJSON(w, &templateData{Form: form})
			return
		}

		// Authenticate the user when the input is correct. If the credentials do not match, the user will receive a generic error message.
		// A generic error prevents from checking to see if an email address exists in our user database and start hacking.
		err = app.models.Users.Authenticate(user.Name, p.Password)

		if err == data.ErrInvalidCredentials {
			form.Errors.AddError("generic", "Username or Password is incorrect")
			// app.render(w, r, "index.html", &templateData{Form: form})
			app.serveAsJSON(w, &templateData{Form: form})
			return
		} else if err != nil {
			app.serverError(w, err)
			return
		}
		authUser, err := app.models.Users.GetByUserCredentials(p.Username)

		if err != nil {
			app.serverError(w, err)
			return
		}
		// Get the token for the current user who is attempting to log in.
		a, err := r.Cookie("session")

		// fmt.Println("Cookie: ", a)

		// expiration := time.Now().Add(5 * time.Minute)
		// cookie := http.Cookie{Name: "newsession", Value: "abcd", Expires: expiration}
		// http.SetCookie(w, &cookie)

		// currentUser := app.contextGetUser(r)

		if err != nil {
			app.serverError(w, err)
		}

		// Add the current cookie (token) to the user's profile in database.
		err = app.models.Users.UpdateByToken(a.Value, authUser.Name)

		if err != nil {
			app.serverError(w, err)
			return
		}

		users, err := app.models.Users.GetAllUsers(authUser.Name)

		if err != nil {
			app.serverError(w, err)
		}

		app.serveAsJSON(w, &templateData{Form: form, AuthenticatedUser: authUser, User: authUser, Users: users})

	}

	// After login redirect the user to the homepage.
	// http.Redirect(w, r, back, http.StatusSeeOther)
}

func (app *application) register(w http.ResponseWriter, r *http.Request) {

	// if r.Header.Get("referer") != "http://localhost:8090/login" && r.Header.Get("referer") != "http://localhost:8090/signup" {
	// 	back = r.Header.Get("referer")
	// }

	switch r.Method {
	case "GET":
		app.home(w, r)
		// app.render(w, r, "register.page.tmpl", &templateData{Form: forms.New(nil)})
		return
	case "POST":

		var regForm forms.RegisterForm

		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&regForm)
		if err != nil {
			app.serverError(w, err)
			return
		}

		// err := r.ParseForm()
		// if err != nil {
		// 	app.serverError(w, err)
		// 	return
		// }

		form := forms.New(r.PostForm)
		v := forms.NewValidator()
		form.Errors = v

		age, err := strconv.Atoi(regForm.Age)
		gender := CheckGender(regForm.Gender)

		user := &data.User{
			Name:    regForm.Username,
			Email:   regForm.Email,
			Forname: regForm.Forname,
			Surname: regForm.Surname,
			Age:     age,
			Gender:  gender,
		}

		err = user.Password.Set(regForm.Password)

		if err != nil {
			app.serverError(w, err)
			return
		}

		// Confirm password match
		v.Check(regForm.Password == regForm.ConfirmPassword, "password", "Passwords don't match")

		// // Validate the user struct and return the error messages to the client if any of
		// // the checks fail.
		data.ValidateUser(v, user)

		if !v.Valid() {
			app.serveAsJSON(w, &templateData{Form: form})
			return
		}

		// // Get the token for the current user who is attempting to register.

		a, err := r.Cookie("session")

		if err != nil {
			app.serverError(w, err)
		}

		//NEW DB FIELDS TO DB
		err = app.models.Users.Insert(user, a.Value)

		if err != nil {
			switch err {
			case data.ErrDuplicateUsername:
				form.Errors.AddError("username", "Username is already in use")
				app.serveAsJSON(w, &templateData{Form: form})
				return
			case data.ErrDuplicateEmail:
				form.Errors.AddError("email", "Email is already in use")
				app.serveAsJSON(w, &templateData{Form: form})
				return
			default:
				app.serverError(w, err)
				return
			}
		}

		// http.Redirect(w, r, back, http.StatusSeeOther)
		authUser := app.contextGetUser(r)

		app.serveAsJSON(w, &templateData{Form: form, AuthenticatedUser: authUser})
		return
	}
}

func (app *application) logout(w http.ResponseWriter, r *http.Request) {
	c := &http.Cookie{Name: "session", Value: uuid.NewV4().String(), Expires: time.Now(), MaxAge: -1}
	http.SetCookie(w, c)
	authUser := app.contextGetUser(r)
	removeSocketReader(authUser.Name)
	// http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) profile(w http.ResponseWriter, r *http.Request) {
	// url := r.URL.Path[9:]

	user := app.contextGetUser(r)
	posts, err := app.models.Posts.GetUserPosts(user.Name)

	if err != nil {
		app.serverError(w, err)
	}

	comments, err := app.models.Comments.GetUserComments(user.Name)
	if err != nil {
		app.serverError(w, err)
	}

	users, err := app.models.Users.GetAllUsers(user.Name)

	if err != nil {
		app.serverError(w, err)
	}

	app.serveAsJSON(w, &templateData{AuthenticatedUser: user, Posts: posts, Comments: comments, Users: users})

	// switch url {
	// case "likes":
	// 	user := app.contextGetUser(r)
	// 	posts, err := app.models.Posts.GetUserLiked(user.Name)
	// 	if err != nil {
	// 		app.serverError(w, err)
	// 	}
	// 	app.render(w, r, "profile.page.tmpl", &templateData{Posts: posts})
	// case "comments":
	// 	user := app.contextGetUser(r)
	// 	comments, err := app.models.Comments.GetUserComments(user.Name)
	// 	if err != nil {
	// 		app.serverError(w, err)
	// 	}
	// 	app.render(w, r, "profile.page.tmpl", &templateData{Comments: comments})
	// default:
	// 	user := app.contextGetUser(r)
	// 	posts, err := app.models.Posts.GetUserPosts(user.Name)
	// 	if err != nil {
	// 		app.serverError(w, err)
	// 	}
	// 	app.render(w, r, "profile.page.tmpl", &templateData{Posts: posts})
	// }
}

func CheckGender(gender string) data.Gender {

	switch gender {
	case "0":
		return data.Male
	case "1":
		return data.Female
	case "2":
		return data.NonBinary
	default:
		return data.Undefined
	}
}
