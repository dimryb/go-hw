package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadDir_ValidVariables(t *testing.T) {
	tests := []struct {
		name     string
		files    map[string]string
		expected Environment
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
		},
		{
			name: "Trimmed var",
			files: map[string]string{
				"TRIMMED_VAR": "   trimmed   ",
			},
			expected: Environment{
				"TRIMMED_VAR": {
					Value:      "   trimmed",
					NeedRemove: false,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()

			for name, content := range tt.files {
				filePath := filepath.Join(tempDir, name)
				err := os.WriteFile(filePath, []byte(content), 0o644)
				require.NoError(t, err, "failed to create test file: %s", filePath)
			}

			env, err := ReadDir(tempDir)
			require.NoError(t, err, "unexpected error: %v", err)
			assert.Equal(t, tt.expected, env, "environment mismatch")
		})
	}
}

func TestReadDir_EmptyAndInvalidVariables(t *testing.T) {
	tests := []struct {
		name     string
		files    map[string]string
		expected Environment
	}{
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
		},
		{
			name: "Invalid var",
			files: map[string]string{
				"INVALID=VAR": "value",
			},
			expected: Environment{},
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
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()

			for name, content := range tt.files {
				filePath := filepath.Join(tempDir, name)
				err := os.WriteFile(filePath, []byte(content), 0o644)
				require.NoError(t, err, "failed to create test file: %s", filePath)
			}

			env, err := ReadDir(tempDir)
			require.NoError(t, err, "unexpected error: %v", err)
			assert.Equal(t, tt.expected, env, "environment mismatch")
		})
	}
}

func TestReadDir_SpecialCharacters(t *testing.T) {
	tests := []struct {
		name     string
		files    map[string]string
		expected Environment
	}{
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
		},
		{
			name: "Multiline file",
			files: map[string]string{
				"MULTILINE_VAR": "first line\nsecond line\nthird line",
			},
			expected: Environment{
				"MULTILINE_VAR": {
					Value:      "first line",
					NeedRemove: false,
				},
			},
		},
		{
			name: "File with spaces only",
			files: map[string]string{
				"SPACES_ONLY": "   ",
			},
			expected: Environment{
				"SPACES_ONLY": {
					Value:      "",
					NeedRemove: true,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()

			for name, content := range tt.files {
				filePath := filepath.Join(tempDir, name)
				err := os.WriteFile(filePath, []byte(content), 0o644)
				require.NoError(t, err, "failed to create test file: %s", filePath)
			}

			env, err := ReadDir(tempDir)
			require.NoError(t, err, "unexpected error: %v", err)
			assert.Equal(t, tt.expected, env, "environment mismatch")
		})
	}
}

func TestReadDir_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		dir      string
		expected Environment
		hasError bool
	}{
		{
			name:     "Non-existent directory",
			dir:      "/non/existent/dir",
			expected: nil,
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env, err := ReadDir(tt.dir)

			if tt.hasError {
				assert.Error(t, err, "expected an error but got none")
				return
			}
			assert.NoError(t, err, "unexpected error: %v", err)
			assert.Equal(t, tt.expected, env, "environment mismatch")
		})
	}
}
