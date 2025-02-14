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
			switch {
			case r == '\\':
				state = escaping
			case unicode.IsDigit(r):
				return "", ErrInvalidString
			default:
				lastSymbol = r
				state = expectAny
			}
		case expectAny:
			switch {
			case r == '\\':
				builder.WriteRune(lastSymbol)
				state = escaping
			case unicode.IsDigit(r):
				count, _ := strconv.Atoi(string(r))
				builder.WriteString(strings.Repeat(string(lastSymbol), count))
				state = expectSymbol
			default:
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
	switch state {
	case expectAny:
		builder.WriteRune(lastSymbol)
	case escaping:
		return "", ErrInvalidString
	}

	return builder.String(), nil
}
