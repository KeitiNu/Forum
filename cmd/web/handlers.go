package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"git.01.kood.tech/roosarula/forum/pkg/data"
	"git.01.kood.tech/roosarula/forum/pkg/forms"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

// Create a variable to store the page where the client was before action (ex. logging in and returning directly to the post)
var back string

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type socketReader struct {
	con  *websocket.Conn
	mode int
	name string
}

var savedsocketreader []*socketReader

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	fmt.Println("In home router")
	switch r.Method {
	case "GET":
		fmt.Println(r.URL.Path)
		// If a user types in an address that doesn't exist, a 404 Error is displayed.
		// if r.URL.Path != "/" {
		// 	w.WriteHeader(404)
		// 	app.render(w, r, "400.page.tmpl", nil)
		// 	return
		// }
		// If the user connects to the home address, the frontpage is displayed.
		categories, err := app.models.Categories.Latest()
		if err != nil {
			app.serverError(w, err)
		}
		// // app.render(w, r, "home.page.tmpl", &templateData{Categories: categories})

		app.render(w, r, "index.html", &templateData{Categories: categories})

		// t, err := template.ParseFiles("ui/html/index.html")
		// if err != nil {
		// 	http.Error(w, "500: internal server error", http.StatusInternalServerError)
		// 	return
		// }

		// j, err := json.Marshal(categories)
		// if err != nil {
		// 	app.serverError(w, err)
		// }

		// t.Execute(w, j)

	case "POST":
		app.serverError(w, errors.New("POST METHOD NOT ALLOWED"))
	}
}

func (app *application) data(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Getting Data, url: ", r.URL.Path)

	switch r.Method {
	case "POST":
		path := r.URL.Path[6:]
		paths := strings.Split(path, "/")
		switch paths[0] {
		case "category":

			app.showCategory(w, r, paths[1])

		case "post":

			app.showPost(w, r, paths[1])
		case "login":

			app.login(w, r)
		case "submit":
			app.submitPost(w, r)

		default:

			categories, _ := app.models.Categories.Latest()

			app.serveAsJSON(w, &templateData{Categories: categories})

		}

	case "GET":
		app.serverError(w, errors.New("POST METHOD NOT ALLOWED"))
	}

}

func (app *application) socket(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":

		if savedsocketreader == nil {
			savedsocketreader = make([]*socketReader, 0)
		}

		defer func() {
			err := recover()
			if err != nil {
				log.Println(err)
			}
			r.Body.Close()

		}()

		con, _ := upgrader.Upgrade(w, r, nil)

		ptrSocketReader := &socketReader{
			con: con,
		}

		savedsocketreader = append(savedsocketreader, ptrSocketReader)

		// ptrSocketReader.con.WriteMessage(websocket.TextMessage, []byte("Greetings from golang"));

		_, message, _ := ptrSocketReader.con.ReadMessage()
		fmt.Println("Message retrieved: ", string(message))

	case "POST":
		app.serverError(w, errors.New("POST METHOD NOT ALLOWED"))
	}
}

func (app *application) login(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET": // When a person clicks the login link, the form appears.
		app.home(w, r)

		// app.render(w, r, "login.page.tmpl", &templateData{Form: forms.New(nil)})
		// if r.Header.Get("referer") != "http://localhost:8090/login" && r.Header.Get("referer") != "http://localhost:8090/signup" {
		// 	back = r.Header.Get("referer")
		// }
		// return
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
		err = app.models.Users.Authenticate(user.Name, form.Get("password"))
		if err == data.ErrInvalidCredentials {
			form.Errors.AddError("generic", "Username or Password is incorrect")
			// app.render(w, r, "index.html", &templateData{Form: form})
			app.serveAsJSON(w, &templateData{Form: form})

			return
		} else if err != nil {
			app.serverError(w, err)
			return
		}

		// Get the token for the current user who is attempting to log in.
		a, err := r.Cookie("session")
		if err != nil {
			app.serverError(w, err)
		}

		// Add the current cookie (token) to the user's profile in database.
		err = app.models.Users.UpdateByToken(a.Value, user.Name)
		if err != nil {
			app.serverError(w, err)
			return
		}
	}
	// After login redirect the user to the homepage.
	http.Redirect(w, r, back, http.StatusSeeOther)
}

func (app *application) register(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("referer") != "http://localhost:8090/login" && r.Header.Get("referer") != "http://localhost:8090/signup" {
		back = r.Header.Get("referer")
	}
	switch r.Method {
	case "GET":
		app.render(w, r, "register.page.tmpl", &templateData{Form: forms.New(nil)})
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

		user := &data.User{
			Name:  form.Get("username"),
			Email: form.Get("email"),
		}
		err = user.Password.Set(form.Get("password"))
		if err != nil {
			app.serverError(w, err)
			return
		}
		// Confirm password match
		plainPass := form.Get("password")
		confirmPass := form.Get("confirm_password")
		v.Check(plainPass == confirmPass, "password", "Passwords don't match")

		// Validate the user struct and return the error messages to the client if any of
		// the checks fail.
		data.ValidateUser(v, user)
		if !v.Valid() {
			app.render(w, r, "register.page.tmpl", &templateData{Form: form})
			return
		}

		// Get the token for the current user who is attempting to register.
		a, err := r.Cookie("session")
		if err != nil {
			app.serverError(w, err)
		}

		err = app.models.Users.Insert(user, a.Value)
		if err != nil {
			switch err {
			case data.ErrDuplicateUsername:
				form.Errors.AddError("username", "Username is already in use")
				app.render(w, r, "register.page.tmpl", &templateData{Form: form})
				return
			case data.ErrDuplicateEmail:
				form.Errors.AddError("email", "Email is already in use")
				app.render(w, r, "register.page.tmpl", &templateData{Form: form})
				return
			default:
				app.serverError(w, err)
				return
			}
		}
		http.Redirect(w, r, back, http.StatusSeeOther)
		return
	}
}

func (app *application) logout(w http.ResponseWriter, r *http.Request) {
	c := &http.Cookie{Name: "session", Value: uuid.NewV4().String(), Expires: time.Now(), MaxAge: -1}
	http.SetCookie(w, c)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) profile(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path[9:]
	switch url {
	case "likes":
		user := app.contextGetUser(r)
		posts, err := app.models.Posts.GetUserLiked(user.Name)
		if err != nil {
			app.serverError(w, err)
		}
		app.render(w, r, "profile.page.tmpl", &templateData{Posts: posts})
	case "comments":
		user := app.contextGetUser(r)
		comments, err := app.models.Comments.GetUserComments(user.Name)
		if err != nil {
			app.serverError(w, err)
		}
		app.render(w, r, "profile.page.tmpl", &templateData{Comments: comments})
	default:
		user := app.contextGetUser(r)
		posts, err := app.models.Posts.GetUserPosts(user.Name)
		if err != nil {
			app.serverError(w, err)
		}
		app.render(w, r, "profile.page.tmpl", &templateData{Posts: posts})
	}
}
