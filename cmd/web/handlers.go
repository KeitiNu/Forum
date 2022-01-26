package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"git.01.kood.tech/roosarula/forum/pkg/data"
	"git.01.kood.tech/roosarula/forum/pkg/forms"
	uuid "github.com/satori/go.uuid"
)

// Create a variable to store the page where the client was before action (ex. logging in and returning directly to the post)
var back string

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":

		// If a user types in an address that doesn't exist, a 404 Error is displayed.
		if r.URL.Path != "/" {
			w.WriteHeader(404)
			app.render(w, r, "400.page.tmpl", nil)
			return
		}
		// If the user connects to the home address, the frontpage is displayed.
		categories, err := app.models.Categories.Latest()
		if err != nil {
			app.serverError(w, err)
		}
		app.render(w, r, "home.page.tmpl", &templateData{Categories: categories})
	case "POST":
		app.serverError(w, errors.New("POST METHOD NOT ALLOWED"))
	}
}

func (app *application) login(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET": // When a person clicks the login link, the form appears.
		app.render(w, r, "login.page.tmpl", &templateData{Form: forms.New(nil)})
		if r.Header.Get("referer") != "http://localhost:8090/login" && r.Header.Get("referer") != "http://localhost:8090/signup" {
			back = r.Header.Get("referer")
		}
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

		// User object
		user := &data.User{
			Name: form.Get("username"),
		}
		err = user.Password.Set(form.Get("password"))
		if err != nil {
			app.serverError(w, err)
		}

		// Validate the user struct and return the error messages to the client if any of
		// the checks fail.
		if data.ValidateLogin(v, user); !v.Valid() {
			app.render(w, r, "login.page.tmpl", &templateData{Form: form})
			return
		}

		// Authenticate the user when the input is correct. If the credentials do not match, the user will receive a generic error message.
		// A generic error prevents from checking to see if an email address exists in our user database and start hacking.
		err = app.models.Users.Authenticate(user.Name, form.Get("password"))
		if err == data.ErrInvalidCredentials {
			form.Errors.AddError("generic", "Username or Password is incorrect")
			app.render(w, r, "login.page.tmpl", &templateData{Form: form})
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

func (app *application) github(w http.ResponseWriter, r *http.Request) {
	// Github returns code in url
	keys, ok := r.URL.Query()["code"]

	if !ok || len(keys[0]) < 1 {
		app.serverError(w, errors.New("paramater 'code' not in url"))
	}
	// Preeparing a request to github to exhange code for user access_token
	code := keys[0]
	client_id := "52144f36461b8f17cc05"
	client_secret := "3230a1d333760f60a8055bf07acd991c4f7882e6"

	postBody, _ := json.Marshal(map[string]string{
		"client_id":     client_id,
		"client_secret": client_secret,
		"code":          code,
	})

	req, err := http.NewRequest("POST", "https://github.com/login/oauth/access_token", bytes.NewBuffer(postBody))
	if err != nil {
		app.serverError(w, err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	// Ecexuting request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		app.serverError(w, err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	resp.Body.Close()
	// Checking if request to github was successful
	if resp.StatusCode != 200 {
		app.serverError(w, errors.New("github API didn't return status 200"))
	}

	var githubResp githubResponse
	err = json.Unmarshal(body, &githubResp)
	if err != nil {
		app.serverError(w, err)
	}

	// Preparing a request to github user API to recieve email address
	req, err = http.NewRequest("GET", "https://api.github.com/user/emails", nil)
	if err != nil {
		app.serverError(w, err)
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Authorization", fmt.Sprintf("token %s", githubResp.Token))
	resp, err = client.Do(req)
	if err != nil {
		app.serverError(w, err)
	}
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		app.serverError(w, err)
	}
	resp.Body.Close()

	var emails []email
	err = json.Unmarshal(body, &emails)
	if err != nil {
		app.serverError(w, err)
	}

	var githubEmail string
	for _, email := range emails {
		if email.Primary {
			githubEmail = email.Email
		}
	}
	// Check if email is already in database
	emailExists, username, err := app.models.Users.EmailExist(githubEmail)
	if err != nil {
		app.serverError(w, err)
	}

	if emailExists {
		// Get the token for the current user who is attempting to log in.
		a, err := r.Cookie("session")
		if err != nil {
			app.serverError(w, err)
		}

		// Add the current cookie (token) to the user's profile in database.
		err = app.models.Users.UpdateByToken(a.Value, username)
		if err != nil {
			app.serverError(w, err)
			return
		}
		// After login redirect the user to the homepage.
		http.Redirect(w, r, back, http.StatusSeeOther)
	}

	t := fmt.Sprintf("email=%s", githubEmail)
	v, err := url.ParseQuery(t)
	if err != nil {
		panic(err)
	}
	app.render(w, r, "github.page.tmpl", &templateData{Form: forms.New(v)})
}

type githubResponse struct {
	Token string `json:"access_token"`
}

type email struct {
	Email    string `json:"email"`
	Primary  bool   `json:"primary"`
	Verified bool   `json:"verified"`
}

func (app *application) registerGithub(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("referer") != "http://localhost:8090/login" && r.Header.Get("referer") != "http://localhost:8090/signup" && r.Header.Get("referer") != "http://localhost:8090/github" {
		back = r.Header.Get("referer")
	}
	switch r.Method {
	case "GET":
		app.render(w, r, "github.page.tmpl", &templateData{Form: forms.New(nil)})
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
			app.render(w, r, "github.page.tmpl", &templateData{Form: form})
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
				app.render(w, r, "github.page.tmpl", &templateData{Form: form})
				return
			default:
				app.serverError(w, err)
				return
			}
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
}
