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

	var errors ValidationErrors
	t := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := t.Field(i)
		value := val.Field(i)

		if field.PkgPath != "" {
			fmt.Println("PkgPath:", field.PkgPath)
			continue
		}

		tag := field.Tag.Get("validate")
		if tag == "" {
			continue
		}

		fieldName := field.Name
		err := validateField(fieldName, value, tag)
		if err != nil {
			errors = append(errors, ValidationError{fieldName, err})
		}
	}

	if len(errors) > 0 {
		return errors
	}

	return nil
}

func validateField(fieldName string, value reflect.Value, tag string) error {
	fmt.Println("fieldName:", fieldName)
	fmt.Println("value:", value)
	fmt.Println("tag:", tag)
	return nil
}
