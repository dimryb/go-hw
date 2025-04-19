package hw09structvalidator

import (
	"fmt"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var msgs []string
	for _, ve := range v {
		msgs = append(msgs, fmt.Sprintf("%s: %v", ve.Field, ve.Err))
	}
	return strings.Join(msgs, "; ")
}

func Validate(v interface{}) error {
	// Place your code here.
	return nil
}
