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

	_ = envDir
	_ = command
	_ = args

	env, err := ReadDir(envDir)
	if err != nil {
		fmt.Printf("Error: failed to read environment directory '%s': %s\n", envDir, err.Error())
	}

	_ = env

	fmt.Println(
		"HELLO is (\"hello\")\n" +
			"BAR is (bar)\n" +
			"FOO is (   foo\nwith new line)\n" +
			"UNSET is ()\n" +
			"ADDED is (from original env)\n" +
			"EMPTY is ()\n" +
			"arguments are arg1=1 arg2=2",
	)
}
