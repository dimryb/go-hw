package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var (
	ErrorValueMustBeStruct    = errors.New("value must be a struct")
	ErrorInvalidRuleFormat    = errors.New("invalid rule format")
	ErrorUnsupportedType      = errors.New("unsupported type")
	ErrorInvalidLenValue      = errors.New("invalid len value")
	ErrorLengthMustBe         = errors.New("length must be")
	ErrorUnknownRuleForString = errors.New("unknown rule for string")
	ErrorValueMustBeOneOf     = errors.New("value must be one of")
	ErrorInvalidRegexp        = errors.New("invalid regexp")
	ErrorDoesNotMatchRegexp   = errors.New("does not match regexp")
	ErrorMinValue             = errors.New("value must be >=")
	ErrorMaxValue             = errors.New("value must be <=")
	ErrorUnknownRuleForInt    = errors.New("unknown rule for int")
)

type Numbers interface {
	int | int8 | int16 | int32 | int64 |
		uint | uint8 | uint16 | uint32 | uint64 | uintptr |
		float32 | float64
}

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	msgs := make([]string, 0, len(v))
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
	rules := strings.Split(tag, "|")
	for _, rule := range rules {
		applyRule(value, rule, fieldName, errs)
	}
}

func applyRule(value reflect.Value, rule string, fieldName string, errs *ValidationErrors) {
	parts := strings.SplitN(rule, ":", 2)
	ruleName := rule
	var ruleValue string

	if len(parts) == 2 {
		ruleName, ruleValue = parts[0], parts[1]
	} else if ruleName != "nested" {
		*errs = append(*errs, ValidationError{
			Field: fieldName,
			Err:   fmt.Errorf("%w: %s", ErrorInvalidRuleFormat, rule),
		})
		return
	}

	switch value.Kind() {
	case reflect.String:
		err := validateString(value.String(), ruleName, ruleValue)
		if err != nil {
			*errs = append(*errs, ValidationError{Field: fieldName, Err: err})
		}
	case reflect.Slice:
		validateSlice(value, ruleName, ruleValue, fieldName, errs)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		num := value.Int()
		err := validateNumber(num, ruleName, ruleValue)
		if err != nil {
			*errs = append(*errs, ValidationError{Field: fieldName, Err: err})
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		num := value.Uint()
		err := validateNumber(num, ruleName, ruleValue)
		if err != nil {
			*errs = append(*errs, ValidationError{Field: fieldName, Err: err})
		}
	case reflect.Float32, reflect.Float64:
		num := value.Float()
		err := validateNumber(num, ruleName, ruleValue)
		if err != nil {
			*errs = append(*errs, ValidationError{Field: fieldName, Err: err})
		}
	case reflect.Struct:
		if ruleName == "nested" {
			nestedErr := Validate(value.Interface())
			if nestedErr != nil {
				var nestedValidationErrs ValidationErrors
				if errors.As(nestedErr, &nestedValidationErrs) {
					for _, ve := range nestedValidationErrs {
						*errs = append(*errs, ValidationError{
							Field: fieldName + "." + ve.Field,
							Err:   ve.Err,
						})
					}
				} else {
					*errs = append(*errs, ValidationError{
						Field: fieldName,
						Err:   nestedErr,
					})
				}
			}
		}
	case reflect.Invalid,
		reflect.Bool,
		reflect.Complex64, reflect.Complex128,
		reflect.Array,
		reflect.Chan,
		reflect.Func,
		reflect.Interface,
		reflect.Map,
		reflect.Ptr,
		reflect.UnsafePointer:
		*errs = append(*errs, ValidationError{
			Field: fieldName,
			Err:   fmt.Errorf("%w: %s", ErrorUnsupportedType, value.Kind()),
		})
	default:
		*errs = append(*errs, ValidationError{
			Field: fieldName,
			Err:   fmt.Errorf("%w: %s", ErrorUnsupportedType, value.Kind()),
		})
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
	case "regexp":
		matched, err := regexp.MatchString(ruleValue, s)
		if err != nil {
			return fmt.Errorf("%w: %s", ErrorInvalidRegexp, ruleValue)
		}
		if !matched {
			return fmt.Errorf("%w: %s", ErrorDoesNotMatchRegexp, ruleValue)
		}
	case "in":
		values := strings.Split(ruleValue, ",")
		if !contains(values, s) {
			return fmt.Errorf("%w %s", ErrorValueMustBeOneOf, ruleValue)
		}
	default:
		return fmt.Errorf("%w: %s", ErrorUnknownRuleForString, ruleName)
	}
	return nil
}

func contains(arr []string, str string) bool {
	for _, v := range arr {
		if v == str {
			return true
		}
	}
	return false
}

func validateNumber[T Numbers](n T, ruleName string, ruleValue string) error {
	switch ruleName {
	case "min":
		minimum, err := parseNumber[T](ruleValue)
		if err != nil {
			return fmt.Errorf("invalid min value: %s", ruleValue)
		}
		if n < minimum {
			return fmt.Errorf("%w %v", ErrorMinValue, minimum)
		}
	case "max":
		maximum, err := parseNumber[T](ruleValue)
		if err != nil {
			return fmt.Errorf("invalid max value: %s", ruleValue)
		}
		if n > maximum {
			return fmt.Errorf("%w %v", ErrorMaxValue, maximum)
		}
	case "in":
		values := strings.Split(ruleValue, ",")
		if !containsNumber(values, n) {
			return fmt.Errorf("%w %s", ErrorValueMustBeOneOf, ruleValue)
		}
	default:
		return fmt.Errorf("%w: %s", ErrorUnknownRuleForInt, ruleName)
	}
	return nil
}

func containsNumber[T Numbers](arr []string, n T) bool {
	for _, v := range arr {
		num, err := parseNumber[T](v)
		if err == nil && num == n {
			return true
		}
	}
	return false
}

func parseNumber[T Numbers](s string) (T, error) {
	var zero T
	switch any(zero).(type) {
	case int, int8, int16, int32, int64:
		val, err := strconv.ParseInt(s, 10, 64)
		return T(val), err
	case uint, uint8, uint16, uint32, uint64, uintptr:
		val, err := strconv.ParseUint(s, 10, 64)
		return T(val), err
	case float32, float64:
		val, err := strconv.ParseFloat(s, 64)
		return T(val), err
	default:
		return zero, fmt.Errorf("unsupported type")
	}
}

func validateSlice(slice reflect.Value, ruleName, ruleValue string, fieldName string, errs *ValidationErrors) {
	for i := 0; i < slice.Len(); i++ {
		elementFieldName := fmt.Sprintf("%s[%d]", fieldName, i)
		applyRule(slice.Index(i), ruleName+":"+ruleValue, elementFieldName, errs)
	}
}
