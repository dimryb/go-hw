package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

var (
	ErrorValueMustBeStruct = errors.New("value must be a struct")
	ErrorInvalidRuleFormat = errors.New("invalid rule format")
	ErrorUnsupportedType   = errors.New("unsupported type")
	ErrorInvalidLenValue   = errors.New("invalid len value")
	ErrorLengthMustBe      = errors.New("length must be")
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
		return ValidationErrors{
			ValidationError{Field: "", Err: ErrorValueMustBeStruct},
		}
	}

	var errs ValidationErrors
	t := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := t.Field(i)
		value := val.Field(i)

		if field.PkgPath != "" {
			continue
		}

		tag := field.Tag.Get("validate")
		if tag == "" {
			continue
		}

		validateField(field.Name, value, tag, &errs)
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}

func validateField(fieldName string, value reflect.Value, tag string, errs *ValidationErrors) {
	fmt.Println("fieldName:", fieldName)
	fmt.Println("tag:", tag)
	rules := strings.Split(tag, "|")
	for _, rule := range rules {
		err := applyRule(value, rule)
		if err != nil {
			*errs = append(*errs, ValidationError{Field: fieldName, Err: err})
		}
	}
}

func applyRule(value reflect.Value, rule string) error {
	fmt.Println("rule:", rule)
	fmt.Println("value:", value)
	parts := strings.SplitN(rule, ":", 2)
	if len(parts) != 2 {
		return fmt.Errorf("%w: %s", ErrorInvalidRuleFormat, rule)
	}
	ruleName, ruleValue := parts[0], parts[1]
	fmt.Println("ruleName:", ruleName)
	fmt.Println("ruleValue:", ruleValue)

	switch value.Kind() {
	case reflect.String:
		return validateString(value.String(), ruleName, ruleValue)
	case reflect.Int:
		return nil
	case reflect.Slice:
		return nil
	default:
		return fmt.Errorf("%w: %s", ErrorUnsupportedType, value.Kind())
	}
}

func validateString(s string, ruleName, ruleValue string) error {
	switch ruleName {
	case "len":
		length, err := strconv.Atoi(ruleValue)
		if err != nil {
			return fmt.Errorf("%w: %s", ErrorInvalidLenValue, ruleValue)
		}
		if len(s) != length {
			return fmt.Errorf("%w %d", ErrorLengthMustBe, length)
		}
	}
	return nil
}
