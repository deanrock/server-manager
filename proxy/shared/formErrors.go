package shared

import (
)

type FormError struct {
	Name string `json:"name"`
	Error string `json:"error"`
}

type FormErrors struct {
	Errors []FormError `json:"errors"`
}

func (fe *FormErrors) Add(name string, error string) {
	fe.Errors = append(fe.Errors, FormError{
		Name: name,
		Error: error,
	})
}

func (fe *FormErrors) HasErrors() bool {
	if len(fe.Errors) > 0 {
		return true
	}

	return false
}
