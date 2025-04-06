package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory '%s': %w", dir, err)
	}

	env := make(Environment)

	for _, file := range files {
		if strings.Contains(file.Name(), "=") {
			continue
		}

		filePath := filepath.Join(dir, file.Name())
		content, err := os.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read file '%s': %w", filePath, err)
		}

		value := strings.TrimRight(string(content), " ")

		if value == "" {
			env[file.Name()] = EnvValue{
				Value:      "",
				NeedRemove: true,
			}
		} else {
			value = strings.ReplaceAll(value, "\x00", "\n")
			env[file.Name()] = EnvValue{
				Value:      value,
				NeedRemove: false,
			}
		}
	}

	return env, nil
}
