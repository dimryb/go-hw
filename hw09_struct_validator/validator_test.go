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

type (
	Meta struct {
		CreatedAt string `validate:"regexp:^\\d{4}-\\d{2}-\\d{2}$"`
	}

	MetaUser struct {
		Meta `validate:"nested"`
		Name string `validate:"len:10"`
	}
)

func TestValidate(t *testing.T) {
	t.Run("Struct validation", TestValidateStruct)
	t.Run("String validation", TestValidateStrings)
	t.Run("Int validation", TestValidateInts)
	t.Run("Float validation", TestValidateFloatTypes)
	t.Run("Slice validation", TestValidateSlices)
}

func assertValidationErrors(t *testing.T, err error, expectedErr error) {
	t.Helper()
	var validationErr ValidationErrors
	switch {
	case errors.As(err, &validationErr):
		found := false
		for _, ve := range validationErr {
			if errors.Is(ve.Err, expectedErr) {
				found = true
				break
			}
		}
		assert.True(t, found, "expected error to contain %v, got %v", expectedErr, err)
	case expectedErr == nil:
		assert.NoError(t, err)
	default:
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
		{
			name: "valid version",
			in: App{
				Version: "v1.00",
			},
			expectedErr: nil,
		},
		{
			name: "version too short",
			in: App{
				Version: "v1.0",
			},
			expectedErr: ErrorLengthMustBe,
		},
		{
			name: "valid token",
			in: Token{
				Header:    []byte{1, 2, 3},
				Payload:   []byte{4, 5, 6},
				Signature: []byte{7, 8, 9},
			},
			expectedErr: nil,
		},
		{
			name: "valid response",
			in: Response{
				Code: 200,
				Body: "OK",
			},
			expectedErr: nil,
		},
		{
			name: "invalid code",
			in: Response{
				Code: 300,
				Body: "Not Found",
			},
			expectedErr: ErrorValueMustBeOneOf,
		},
		{
			name: "empty body",
			in: Response{
				Code: 404,
				Body: "",
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
		{
			name: "valid int8",
			in: struct {
				Field int8 `validate:"min:10|max:20"`
			}{
				Field: 15,
			},
			expectedErr: nil,
		},
		{
			name: "invalid int8",
			in: struct {
				Field int8 `validate:"min:100"`
			}{
				Field: 50,
			},
			expectedErr: ErrorMinValue,
		},
		{
			name: "valid uint64",
			in: struct {
				Field uint64 `validate:"max:100"`
			}{
				Field: 50,
			},
			expectedErr: nil,
		},
		{
			name: "invalid uint64",
			in: struct {
				Field uint64 `validate:"min:100"`
			}{
				Field: 50,
			},
			expectedErr: ErrorMinValue,
		},
		{
			name: "valid uint8",
			in: struct {
				Field uint8 `validate:"min:10|max:20"`
			}{
				Field: 15,
			},
			expectedErr: nil,
		},
		{
			name: "invalid uint8",
			in: struct {
				Field uint8 `validate:"min:100"`
			}{
				Field: 50,
			},
			expectedErr: ErrorMinValue,
		},
		{
			name: "valid uintptr",
			in: struct {
				Field uintptr `validate:"in:1,2,3"`
			}{
				Field: 2,
			},
			expectedErr: nil,
		},
		{
			name: "invalid uintptr",
			in: struct {
				Field uintptr `validate:"min:100"`
			}{
				Field: 50,
			},
			expectedErr: ErrorMinValue,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.in)
			assertValidationErrors(t, err, tt.expectedErr)
		})
	}
}

func TestValidateFloatTypes(t *testing.T) {
	tests := []struct {
		name        string
		in          interface{}
		expectedErr error
	}{
		{
			name: "valid float32",
			in: struct {
				Field float32 `validate:"min:10.5|max:20.5"`
			}{
				Field: 15.0,
			},
			expectedErr: nil,
		},
		{
			name: "invalid float32",
			in: struct {
				Field float32 `validate:"max:10.0"`
			}{
				Field: 15.0,
			},
			expectedErr: ErrorMaxValue,
		},
		{
			name: "valid float64",
			in: struct {
				Field float64 `validate:"in:1.0,2.0,3.0"`
			}{
				Field: 2.0,
			},
			expectedErr: nil,
		},
		{
			name: "invalid float64",
			in: struct {
				Field float64 `validate:"min:10.0"`
			}{
				Field: 5.0,
			},
			expectedErr: ErrorMinValue,
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
	}{
		{
			name: "slice element length mismatch",
			in: struct {
				Field []string `validate:"len:5"`
			}{
				Field: []string{"short", "longer"},
			},
			expectedErr: ErrorLengthMustBe,
		},
		{
			name: "valid slice",
			in: struct {
				Field []string `validate:"len:5"`
			}{
				Field: []string{"exact", "exact"},
			},
			expectedErr: nil,
		},
		{
			name: "slice element not in list",
			in: struct {
				Field []int `validate:"in:1,2,3"`
			}{
				Field: []int{1, 4, 3},
			},
			expectedErr: ErrorValueMustBeOneOf,
		},
		{
			name: "slice element less than min",
			in: struct {
				Field []int `validate:"min:10"`
			}{
				Field: []int{5, 15, 20},
			},
			expectedErr: ErrorMinValue,
		},
		{
			name: "slice element greater than max",
			in: struct {
				Field []int `validate:"max:10"`
			}{
				Field: []int{5, 15, 3},
			},
			expectedErr: ErrorMaxValue,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.in)
			assertValidationErrors(t, err, tt.expectedErr)
		})
	}
}

func TestValidateNestedStructuresWithNestedTag(t *testing.T) {
	tests := []struct {
		name        string
		in          interface{}
		expectedErr error
	}{
		{
			name: "valid nested structure with nested tag",
			in: MetaUser{
				Meta: Meta{
					CreatedAt: "2023-10-01",
				},
				Name: "JohnDoe123",
			},
			expectedErr: nil,
		},
		{
			name: "invalid nested structure with nested tag",
			in: MetaUser{
				Meta: Meta{
					CreatedAt: "2023/10/01",
				},
				Name: "JohnDoe123",
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
