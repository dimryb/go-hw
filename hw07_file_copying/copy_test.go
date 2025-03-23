package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCopy(t *testing.T) {
	testdataDir := "testdata"
	tests := []struct {
		name       string
		offset     int64
		limit      int64
		outputFile string // Имя эталонного файла в testdata
	}{
		{
			name:       "offset=0, limit=0 (copy all)",
			offset:     0,
			limit:      0,
			outputFile: "out_offset0_limit0.txt",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputFile := filepath.Join(testdataDir, "input.txt")

			tmpFile, err := os.CreateTemp("", "test_copy_*.txt")
			if err != nil {
				t.Fatalf("failed to create temporary file: %v", err)
			}
			defer os.Remove(tmpFile.Name())

			tmpFilePath := tmpFile.Name()
			tmpFile.Close()

			err = Copy(inputFile, tmpFilePath, tt.offset, tt.limit)
			assert.NoError(t, err, "Copy should not return an error")

			result, err := os.ReadFile(tmpFilePath)
			assert.NoError(t, err, "failed to read temporary file")

			expectedFile := filepath.Join(testdataDir, tt.outputFile)
			expected, err := os.ReadFile(expectedFile)
			assert.NoError(t, err, "failed to read expected file")

			if string(result) != string(expected) {
				t.Errorf("unexpected output:\ngot:\n%s\nexpected:\n%s", result, expected)
			}

			assert.Equal(t, string(expected), string(result),
				"output content does not match expected content")
		})
	}
}
