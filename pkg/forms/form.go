package forms

import (
	"net/url"
)

// Anonymously embeds a url.Values object
// (to hold the form data) and an Errors field to hold any validation errors
// for the form data.
type Form struct {
	url.Values
	Errors *Validator
}

type LoginForm struct {
	Username string
	Password string
}

func New(data url.Values) *Form {
	return &Form{
		data,
		&Validator{},
	}
}
