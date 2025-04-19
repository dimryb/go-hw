package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

var (
	ErrorValidateValueMustBeStruct = errors.New("value must be a struct")
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
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Struct {
		return ErrorValidateValueMustBeStruct
	}

	return nil
}
