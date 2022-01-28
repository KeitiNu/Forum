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
	"strings"

	"git.01.kood.tech/roosarula/forum/pkg/data"
	"git.01.kood.tech/roosarula/forum/pkg/forms"
)

type githubResponse struct {
	Token string `json:"access_token"`
}

type email struct {
	Email    string `json:"email"`
	Primary  bool   `json:"primary"`
	Verified bool   `json:"verified"`
}

type googleData struct {
	Email string `json:"email"`
}

// var (
// 	googleOauthConfig = &oauth2.Config{
// 		RedirectURL:  "http://localhost:8090/google",
// 		ClientID:     "1087259911821-nmcttkbvat9mmkjqrrjl16nahcofdts0.apps.googleusercontent.com",
// 		ClientSecret: "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",
// 		Scopes: []string{"https://www.googleapis.com/auth/userinfo.profile",
// 			"https://www.googleapis.com/auth/userinfo.email"},
// 		Endpoint: google.Endpoint,
// 	}
// 	oauthStateString = "random"
// )

// google is for logging in through google account. Receives accesstoken from google and makes request to google API to receive user email
func (app *application) google(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")

	// Preparing a request to github to exhange code for user access_token
	client_id := "1087259911821-nmcttkbvat9mmkjqrrjl16nahcofdts0.apps.googleusercontent.com"
	client_secret := "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"

	params := url.Values{}
	params.Add("client_id", client_id)
	params.Add("client_secret", client_secret)
	params.Add("grant_type", "authorization_code")
	params.Add("code", code)
	params.Add("redirect_uri", "http://localhost:8090/google")

	// postBody, _ := json.Marshal(map[string]string{
	// 	"grant_type": "authorization_code",
	// 	"client_id":     client_id,
	// 	"client_secret": client_secret,
	// 	"code":          code,
	// 	"redirect_uri": "http://localhost:8090/google"
	// })

	req, err := http.NewRequest("POST", "https://oauth2.googleapis.com/token", strings.NewReader(params.Encode()))
	if err != nil {
		app.serverError(w, err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

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

	var githubResp githubResponse
	err = json.Unmarshal(body, &githubResp)
	if err != nil {
		app.serverError(w, err)
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + githubResp.Token)
	if err != nil {
		app.serverError(w, err)
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		app.serverError(w, err)
	}

	var data googleData
	err = json.Unmarshal(contents, &data)
	if err != nil {
		app.serverError(w, err)
	}

	// Check if email is already in database
	emailExists, username, err := app.models.Users.EmailExist(data.Email)
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
		}
		// After login redirect the user to the homepage.
		http.Redirect(w, r, back, http.StatusSeeOther)
	}

	t := fmt.Sprintf("email=%s", data.Email)
	v, err := url.ParseQuery(t)
	if err != nil {
		app.serverError(w, err)
	}
	app.render(w, r, "otherRegister.page.tmpl", &templateData{Form: forms.New(v)})
}

// github  is for logging in through github. Makes first request to github.com to recieve client token, and second request to github API to get client email
func (app *application) github(w http.ResponseWriter, r *http.Request) {
	// Github returns code in url
	keys, ok := r.URL.Query()["code"]

	if !ok || len(keys[0]) < 1 {
		app.serverError(w, errors.New("paramater 'code' not in url"))
	}
	// Preparing a request to github to exhange code for user access_token
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
	app.render(w, r, "otherRegister.page.tmpl", &templateData{Form: forms.New(v)})
}

// registerThirdParty registers a new account with the email recieved either from github or google.
func (app *application) registerThirdParty(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		app.render(w, r, "otherRegister.page.tmpl", &templateData{Form: forms.New(nil)})
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
			app.render(w, r, "otherRegister.page.tmpl", &templateData{Form: form})
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
				app.render(w, r, "otherRegister.page.tmpl", &templateData{Form: form})
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

// // handleGoogleLogin redirects to google's login site
// func (app *application) handleGoogleLogin(w http.ResponseWriter, r *http.Request) {
// 	http.Redirect(w, r, googleOauthConfig.AuthCodeURL(oauthStateString), http.StatusTemporaryRedirect)
// }
// "https://accounts.google.com/o/oauth2/auth/oauthchooseaccount
// client_id=1087259911821-nmcttkbvat9mmkjqrrjl16nahcofdts0.apps.googleusercontent.com
// redirect_uri=http%3A%2F%2Flocalhost%3A8090%2Fgoogle
// response_type=code
// scope=https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fuserinfo.profile%20https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fuserinfo.email
// state=random
// flowName=GeneralOAuthFlow"

// "https://accounts.google.com/o/oauth2/auth/oauthchooseaccount?client_id=1087259911821-nmcttkbvat9mmkjqrrjl16nahcofdts0.apps.googleusercontent.com&redirect_uri=http%3A%2F%2Flocalhost%3A8090%2Fgoogle&response_type=code&scope=https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fuserinfo.profile%20https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fuserinfo.email"&flowName=GeneralOAuthFlow"