package shared

import (
)

type FormError struct {
	Error string `json:"error"`
}

type FormErrors struct {
	Errors map[string][]FormError `json:"errors"`
}

func NewFormErrors() FormErrors {
	f := FormErrors{
		Errors:  make(map[string][]FormError),
	}

	return f
}

func (fe *FormErrors) Add(name string, error string) {
	e := FormError{
		Error: error,
	}

	if _, ok := fe.Errors[name]; ok {
		fe.Errors[name] = append(fe.Errors[name], e)
	}else{
		fe.Errors[name] = []FormError{e}
	}
}

func (fe *FormErrors) HasErrors() bool {
	if len(fe.Errors) > 0 {
		return true
	}

	return false
}
