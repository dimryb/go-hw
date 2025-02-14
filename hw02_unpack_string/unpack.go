package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	if str == "" {
		return "", nil
	}

	const (
		expectSymbol = iota
		expectAny
		escaping
	)

	state := expectSymbol
	var lastSymbol rune
	var builder strings.Builder

	for _, r := range str {
		switch state {
		case expectSymbol:
			if r == '\\' {
				state = escaping
			} else if unicode.IsDigit(r) {
				return "", ErrInvalidString
			} else {
				lastSymbol = r
				state = expectAny
			}
		case expectAny:
			if r == '\\' {
				builder.WriteRune(lastSymbol)
				state = escaping
			} else if unicode.IsDigit(r) {
				count, _ := strconv.Atoi(string(r))
				builder.WriteString(strings.Repeat(string(lastSymbol), count))
				state = expectSymbol
			} else {
				builder.WriteRune(lastSymbol)
				lastSymbol = r
			}
		case escaping:
			if !unicode.IsDigit(r) && r != '\\' {
				return "", ErrInvalidString
			}
			lastSymbol = r
			state = expectAny
		}
	}
	if state == expectAny {
		builder.WriteRune(lastSymbol)
	} else if state == escaping {
		return "", ErrInvalidString
	}

	return builder.String(), nil
}
