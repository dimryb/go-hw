package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	t.Run("Struct validation", TestValidateStruct)
	t.Run("String validation", TestValidateStrings)
	t.Run("Int validation", TestValidateInts)
	t.Run("Slice validation", TestValidateSlices)
}

func assertValidationErrors(t *testing.T, err error, expectedErr error) {
	var validationErr ValidationErrors
	if errors.As(err, &validationErr) {
		found := false
		for _, ve := range validationErr {
			if errors.Is(ve.Err, expectedErr) {
				found = true
				break
			}
		}
		assert.True(t, found, "expected error to contain %v, got %v", expectedErr, err)
	} else if expectedErr == nil {
		assert.NoError(t, err)
	} else {
		t.Errorf("expected ValidationErrors containing %v, got: %v", expectedErr, err)
	}
}

func TestValidateStruct(t *testing.T) {
	tests := []struct {
		name        string
		in          interface{}
		expectedErr error
	}{
		{
			name:        "non-struct input",
			in:          0,
			expectedErr: ErrorValueMustBeStruct,
		},
		{
			name: "valid user",
			in: User{
				ID:     "123e4567-e89b-12d3-a456-426614174000",
				Age:    19,
				Email:  "valid_user@example.com",
				Role:   "admin",
				Phones: []string{"12345678901"},
			},
			expectedErr: nil,
		},
		{
			name: "invalid rule format",
			in: struct {
				Field string `validate:"min"`
			}{
				Field: "test",
			},
			expectedErr: ErrorInvalidRuleFormat,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.in)
			assertValidationErrors(t, err, tt.expectedErr)
		})
	}
}

func TestValidateStrings(t *testing.T) {
	tests := []struct {
		name        string
		in          interface{}
		expectedErr error
	}{
		{
			name: "invalid len value",
			in: struct {
				Field string `validate:"len:abc"`
			}{
				Field: "test",
			},
			expectedErr: ErrorInvalidLenValue,
		},
		{
			name: "length must be",
			in: struct {
				Field string `validate:"len:5"`
			}{
				Field: "long test string",
			},
			expectedErr: ErrorLengthMustBe,
		},
		{
			name: "unknown rule for string",
			in: struct {
				Field string `validate:"unknown:123"`
			}{
				Field: "test",
			},
			expectedErr: ErrorUnknownRuleForString,
		},
		{
			name: "value must be one of",
			in: struct {
				Field string `validate:"in:admin,user"`
			}{
				Field: "guest",
			},
			expectedErr: ErrorValueMustBeOneOf,
		},
		{
			name: "invalid regexp",
			in: struct {
				Field string `validate:"regexp:[a-z"`
			}{
				Field: "test",
			},
			expectedErr: ErrorInvalidRegexp,
		},
		{
			name: "does not match regexp",
			in: struct {
				Field string `validate:"regexp:^\\d+$"`
			}{
				Field: "abc",
			},
			expectedErr: ErrorDoesNotMatchRegexp,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.in)
			assertValidationErrors(t, err, tt.expectedErr)
		})
	}
}

func TestValidateInts(t *testing.T) {
	tests := []struct {
		name        string
		in          interface{}
		expectedErr error
	}{
		{
			name: "value less than min",
			in: struct {
				Field int `validate:"min:10"`
			}{
				Field: 5,
			},
			expectedErr: ErrorMinValue,
		},
		{
			name: "value greater than max",
			in: struct {
				Field int `validate:"max:10"`
			}{
				Field: 15,
			},
			expectedErr: ErrorMaxValue,
		},
		{
			name: "value not in list",
			in: struct {
				Field int `validate:"in:1,2,3"`
			}{
				Field: 4,
			},
			expectedErr: ErrorValueMustBeOneOf,
		},
		{
			name: "unknown rule for int",
			in: struct {
				Field int `validate:"unknown:10"`
			}{
				Field: 5,
			},
			expectedErr: ErrorUnknownRuleForInt,
		},
		{
			name: "valid int",
			in: struct {
				Field int `validate:"min:10|max:20|in:15,20"`
			}{
				Field: 15,
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.in)
			assertValidationErrors(t, err, tt.expectedErr)
		})
	}
}

func TestValidateSlices(t *testing.T) {
	tests := []struct {
		name        string
		in          interface{}
		expectedErr error
	}{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.in)
			assertValidationErrors(t, err, tt.expectedErr)
		})
	}
}
