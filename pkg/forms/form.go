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
type RegisterForm struct {
	Username        string
	Email           string
	Forname         string
	Surname         string
	Age             string
	Gender          string
	Password        string
	ConfirmPassword string
}

type ChatForm struct {
	Message     string
	RecipientId string
	UserId      string
}

type ChatBoxForm struct {
	User      string
	Recipient string
	Offset    int
}

func New(data url.Values) *Form {
	return &Form{
		data,
		&Validator{},
	}
}
