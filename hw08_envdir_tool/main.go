package main

import (
	"fmt"
)

func main() {
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
