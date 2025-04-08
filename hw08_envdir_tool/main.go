package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go-envdir /path/to/env/dir command [args...]")
		os.Exit(1)
	}

	envDir := os.Args[1]
	command := os.Args[2]
	args := os.Args[3:]

	env, err := ReadDir(envDir)
	if err != nil {
		fmt.Printf("Error: failed to read environment directory '%s': %s\n", envDir, err.Error())
		os.Exit(1)
	}

	code := RunCmd(append([]string{command}, args...), env)
	os.Exit(code)
}
