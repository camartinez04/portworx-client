package main

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/asaskevich/govalidator"
)

// Add adds an error message for a given form field
func (e ErrorsForm) Add(field, message string) {
	e[field] = append(e[field], message)
}

// Get returns the first error message
func (e ErrorsForm) Get(field string) string {
	es := e[field]
	if len(es) == 0 {

		return ""

	}
	return es[0]
}

// Form creates a custom form struct, embeds a url.Values object
type Form struct {
	url.Values
	Errors ErrorsForm
}

func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

// New initializes a form struct
func New(data url.Values) *Form {
	return &Form{
		data,
		ErrorsForm(map[string][]string{}),
	}
}

func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be blank")
		}
	}
}

// Has checks if form field is in post and not empty
func (f *Form) Has(field string) bool {
	x := f.Get(field)
	if x == "" {
		return false
	}
	return true
}

func (f *Form) MinLength(field string, length int) bool {

	x := f.Get(field)
	if len(x) < length {
		f.Errors.Add(field, fmt.Sprintf("This field must be at least %d characters long", length))
		return false
	}

	return true

}

func (f *Form) IsEmail(field string) {
	if !govalidator.IsEmail(f.Get(field)) {
		f.Errors.Add(field, "Invalid email address")
	}
}
