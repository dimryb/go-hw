package main

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
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
		{
			name:       "offset=0, limit=10",
			offset:     0,
			limit:      10,
			outputFile: "out_offset0_limit10.txt",
		},
		{
			name:       "offset=0, limit=1000",
			offset:     0,
			limit:      1000,
			outputFile: "out_offset0_limit1000.txt",
		},
		{
			name:       "offset=0, limit=10000",
			offset:     0,
			limit:      10000,
			outputFile: "out_offset0_limit10000.txt",
		},
		{
			name:       "offset=100, limit=1000",
			offset:     100,
			limit:      1000,
			outputFile: "out_offset100_limit1000.txt",
		},
		{
			name:       "offset=6000, limit=1000",
			offset:     6000,
			limit:      1000,
			outputFile: "out_offset6000_limit1000.txt",
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

			err = Copy(inputFile, tmpFilePath, tt.offset, tt.limit, nil)
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

func TestCopyNonRegularFile(t *testing.T) {
	dir, err := os.MkdirTemp("", "test_dir")
	if err != nil {
		t.Fatalf("failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(dir)

	outputFile := filepath.Join(dir, "output.txt")
	err = Copy(dir, outputFile, 0, 10, nil)
	assert.Error(t, err, "Copy should return an error for non-regular files")
	assert.True(t, errors.Is(err, ErrUnsupportedFile), "Error should be ErrUnsupportedFile")
}

func TestCopyNullDeviceUnsupportedFile(t *testing.T) {
	nullDevice := "/dev/null"
	if runtime.GOOS == "windows" {
		nullDevice = "NUL"
	}

	tmpFile, err := os.CreateTemp("", "test_copy_*.txt")
	if err != nil {
		t.Fatalf("failed to create temporary file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFilePath := tmpFile.Name()
	tmpFile.Close()

	err = Copy(nullDevice, tmpFilePath, 0, 10, nil)
	assert.Error(t, err, "Copy should return an error for unsupported files")
	assert.True(t, errors.Is(err, ErrUnsupportedFile), "Error should be ErrUnsupportedFile")
}
