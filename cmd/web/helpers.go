package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"runtime/debug"
)

// We have error messages and webpage rendering functions here, which are repeated several times and
// would make the code longer if typed everytime. You can see render and error helpers in handlers.go file.

// The serverError helper writes an error message and stack trace to the errorLog,
// then sends a generic 500 Internal Server Error response to the user.
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	err = app.errorLog.Output(2, trace)
	if err != nil {
		app.serverError(w, err)
	}
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// The clientError helper sends a specific status code and corresponding description
// to the user. We'll use this to send responses like 400 "Bad
// Request" when there's a problem with the request that the user sent.
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// This is simply a convenience wrapper around clientError which sends a
// 404 Not Found response to the user.
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) render2(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
	// Retrieve the appropriate template set from the cache based on the page name
	// (like 'home.page.tmpl'). If no entry exists in the cache with the
	// provided name, call the serverError.
	ts, ok := app.templateCache[name]
	if !ok {
		app.serverError(w, fmt.Errorf("the template %s does not exist", name))
		return
	}

	// At first we write template to buffer to know if there would be any errors. If there are no errors we will write it to client.
	buf := new(bytes.Buffer)
	// Execute the template set, passing in any dynamic data.
	// Add default data adds our user information to each template (if user is logged in)
	err := ts.Execute(buf, app.addDefaultData(td, r))
	if err != nil {
		app.serverError(w, err)
	}

	// If there aren't any errors we display it to our client.
	_, err = buf.WriteTo(w)
	if err != nil {
		app.serverError(w, err)
	}
}



func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
	// Retrieve the appropriate template set from the cache based on the page name
	// (like 'home.page.tmpl'). If no entry exists in the cache with the
	// provided name, call the serverError.
	ts, ok := app.templateCache[name]
	if !ok {
		app.serverError(w, fmt.Errorf("the template %s does not exist", name))
		return
	}


	
	// At first we write template to buffer to know if there would be any errors. If there are no errors we will write it to client.
	buf := new(bytes.Buffer)
	// Execute the template set, passing in any dynamic data.
	// Add default data adds our user information to each template (if user is logged in)
	// err := ts.Execute(buf, app.addDefaultData(td, r))

	err := ts.Execute(buf, app.addDefaultData(td, r))
	if err != nil {
		app.serverError(w, err)
	}

	// If there aren't any errors we display it to our client.
	_, err = buf.WriteTo(w)
	if err != nil {
		app.serverError(w, err)
	}
}
// Adding data to every template. For the moment we use it to add information about user to each page.
func (app *application) addDefaultData(td *templateData, r *http.Request) *templateData {
	if td == nil {
		td = &templateData{}
	}
	td.AuthenticatedUser = app.contextGetUser(r)
	if td.AuthenticatedUser == nil {
		return td
	} else {
		td.UserVotes = app.models.Posts.GetUserVotes(td.AuthenticatedUser.Name)
	}
	return td
}


func (app *application) serveAsJSON(w http.ResponseWriter, data *templateData){

	j, err := json.Marshal(data)
	if err != nil {
		app.serverError(w, err)
	}

	io.Copy(w, bytes.NewReader(j))
}
