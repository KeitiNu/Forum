package forms

// Define a new errors type, which we will use to hold the validation error
// messages for forms. The name of the form field will be used as the key in
// this map.
// Example errors["username"]:["Username cannot contain the word SINEP"]
type errors map[string][]string

// Adding the error.
func (e errors) Add(field, message string) {
	e[field] = append(e[field], message)
}

// Getting the error from a map of errors. If there is some data in for example ["email"] we will get a string.
func (e errors) Get(field string) string {
	es := e[field]
	if len(es) == 0 {
		return ""
	}
	return es[0]
}
