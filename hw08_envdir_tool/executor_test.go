package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRunCmd(t *testing.T) {
	tempDir := t.TempDir()

	testScript := createTestScript(t, tempDir)

	tests := []struct {
		name       string
		cmd        []string
		env        Environment
		expected   int
		shouldFail bool
	}{
		{
			name: "Successful command execution",
			cmd:  []string{"bash", testScript, "arg1", "arg2"}, // С "bash" работает на винде тоже
			env: Environment{
				"TEST_VAR": {Value: "test_value", NeedRemove: false},
			},
			expected: 0,
		},
		{
			name:     "Command not found",
			cmd:      []string{"non_existent_command"},
			env:      nil,
			expected: 1,
		},
		{
			name: "Environment variable passed correctly",
			cmd:  []string{"bash", testScript, "printenv", "TEST_VAR"},
			env: Environment{
				"TEST_VAR": {Value: "test_value", NeedRemove: false},
			},
			expected: 0,
		},
		{
			name: "Environment variable removed",
			cmd:  []string{"bash", testScript, "printenv", "REMOVED_VAR"},
			env: Environment{
				"REMOVED_VAR": {Value: "", NeedRemove: true},
			},
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code := RunCmd(tt.cmd, tt.env)
			if code != tt.expected {
				t.Errorf("expected exit code %d, got %d", tt.expected, code)
			}
		})
	}
}

func createTestScript(t *testing.T, dir string) string {
	t.Helper()
	scriptName := "test_script.sh"
	scriptPath := filepath.Join(dir, scriptName)
	scriptContent := `#!/bin/sh
if [ "$1" = "printenv" ]; then
    if [ -z "${!2}" ]; then
        exit 1
    fi
    echo "$2=$3"
    exit 0
fi
exit 0`
	err := os.WriteFile(scriptPath, []byte(scriptContent), 0o755)
	if err != nil {
		t.Fatalf("failed to create test script: %v", err)
	}
	return scriptPath
}
