package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) == 0 {
		fmt.Println("Error: command is empty")
		return 1
	}
	// #nosec G204: Command arguments are trusted
	command := exec.Command(cmd[0], cmd[1:]...)
	command.Env = getEnvironment(env)

	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	err := command.Run()
	if err != nil {
		if code, ok := getExitCode(err); ok {
			return code
		}
		return 1
	}
	return 0
}

func getEnvironment(env Environment) []string {
	if env == nil {
		return nil
	}
	var environment []string

	for _, envVar := range os.Environ() {
		index := strings.Index(envVar, "=")
		if index == -1 {
			continue
		}
		key := envVar[:index]
		if _, exists := env[key]; !exists {
			environment = append(environment, envVar)
		}
	}

	for key, val := range env {
		if !val.NeedRemove {
			environment = append(environment, fmt.Sprintf("%s=%s", key, val.Value))
		}
	}
	return environment
}

func getExitCode(err error) (int, bool) {
	var exitError *exec.ExitError
	if errors.As(err, &exitError) {
		if status, ok := exitError.Sys().(interface{ ExitStatus() int }); ok {
			return status.ExitStatus(), true
		}
	}
	return 0, false
}
