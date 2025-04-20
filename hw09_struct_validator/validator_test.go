package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
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
				ID:     "123e4567-e89b",
				Age:    19,
				Email:  "valid.user@example.com",
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
			name: "unsupported type",
			in: struct {
				Field complex128 `validate:"len:5"`
			}{
				Field: 42 + 2i,
			},
			expectedErr: ErrorUnsupportedType,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)
			if tt.expectedErr == nil {
				assert.NoError(t, err)
			} else {
				fmt.Println("expected:", tt.expectedErr)
				fmt.Println("actual:", err)
				var validationErr ValidationErrors
				if errors.As(err, &validationErr) {
					found := false
					for _, ve := range validationErr {
						fmt.Println("ve:", ve)
						if errors.Is(ve.Err, tt.expectedErr) {
							found = true
							break
						}
					}
					assert.True(t, found, "expected error to contain %v, got %v", tt.expectedErr, err)
				} else {
					assert.Fail(
						t, "expected ValidationErrors, got other error",
						"expected: %v, got: %v",
						tt.expectedErr, err,
					)
				}
				fmt.Println("validationErr:", validationErr)
			}
		})
	}
}
