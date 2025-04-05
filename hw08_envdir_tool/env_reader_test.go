package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	tests := []struct {
		name     string
		files    map[string]string
		dir      string
		expected Environment
		hasError bool
	}{
		{
			name: "Valid var",
			files: map[string]string{
				"VALID_VAR": "value",
			},
			expected: Environment{
				"VALID_VAR": {
					Value:      "value",
					NeedRemove: false,
				},
			},
			hasError: false,
		},
		{
			name: "Empty var",
			files: map[string]string{
				"EMPTY_VAR": "",
			},
			expected: Environment{
				"EMPTY_VAR": {
					Value:      "",
					NeedRemove: true,
				},
			},
			hasError: false,
		},
		{
			name: "With null",
			files: map[string]string{
				"WITH_NULL": "hello\x00world",
			},
			expected: Environment{
				"WITH_NULL": {
					Value:      "hello\nworld",
					NeedRemove: false,
				},
			},
			hasError: false,
		},
		{
			name: "Trimmed var",
			files: map[string]string{
				"TRIMMED_VAR": "   trimmed   ",
			},
			expected: Environment{
				"TRIMMED_VAR": {
					Value:      "trimmed",
					NeedRemove: false,
				},
			},
			hasError: false,
		},
		{
			name: "Invalid var",
			files: map[string]string{
				"INVALID=VAR": "value",
			},
			expected: Environment{},
			hasError: false,
		},
		{
			name: "With mixed files",
			files: map[string]string{
				"VALID_VAR":   "value",
				"EMPTY_VAR":   "",
				"WITH_NULL":   "hello\x00world",
				"INVALID=VAR": "value",
				"TRIMMED_VAR": "   trimmed   ",
			},
			expected: Environment{
				"VALID_VAR": {
					Value:      "value",
					NeedRemove: false,
				},
				"EMPTY_VAR": {
					Value:      "",
					NeedRemove: true,
				},
				"WITH_NULL": {
					Value:      "hello\nworld",
					NeedRemove: false,
				},
				"TRIMMED_VAR": {
					Value:      "trimmed",
					NeedRemove: false,
				},
			},
			hasError: false,
		},
		{
			name:     "Non-existent directory",
			files:    nil,
			dir:      "/non/existent/dir",
			expected: nil,
			hasError: true,
		},
		{
			name: "With no valid files",
			files: map[string]string{
				"INVALID=VAR": "value",
				"EMPTY_VAR":   "",
			},
			expected: Environment{
				"EMPTY_VAR": {
					Value:      "",
					NeedRemove: true,
				},
			},
			hasError: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var tempDir string

			if tt.files != nil {
				tempDir = t.TempDir()

				for name, content := range tt.files {
					filePath := filepath.Join(tempDir, name)
					err := os.WriteFile(filePath, []byte(content), 0o644)
					require.NoError(t, err, "failed to create test file: %s", filePath)
				}
			} else {
				tempDir = tt.dir
			}

			env, err := ReadDir(tempDir)

			if tt.hasError {
				assert.Error(t, err, "expected an error but got none")
				return
			}
			assert.NoError(t, err, "unexpected error: %v", err)
			assert.Equal(t, tt.expected, env, "environment mismatch")
		})
	}
}
