package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(inputStr string) (string, error) {
	var b strings.Builder

	for pos := 0; pos < len(inputStr); {
		substr, err := findNextSubstr(inputStr[pos:])
		if err != nil {
			return "", err
		}
		b.WriteString(unpackSubstr(substr))
		pos += len(substr)
	}
	return b.String(), nil
}

func findNextSubstr(inputStr string) (string, error) {
	runes := []rune(inputStr)
	screenOffset := 0

	if unicode.IsDigit(runes[0]) {
		return "", ErrInvalidString
	}

	if runes[0] == '\\' {
		if len(runes) < 2 {
			return "", ErrInvalidString
		}
		if !unicode.IsDigit(runes[1]) && runes[1] != '\\' {
			return "", ErrInvalidString
		}
		screenOffset = 1
	}

	isSubstrOneChar := len(runes) < 2+screenOffset || !unicode.IsDigit(runes[1+screenOffset])
	isMultipleCyphersInSubstr := len(runes) > 2+screenOffset && unicode.IsDigit(runes[2+screenOffset])

	if isSubstrOneChar {
		return inputStr[:1+screenOffset], nil
	} else {
		if isMultipleCyphersInSubstr {
			return "", ErrInvalidString
		}
		return inputStr[:2+screenOffset], nil
	}
}

func unpackSubstr(substr string) string {
	runes := []rune(substr)
	if runes[0] == '\\' {
		substr = substr[1:]
	}
	if len(substr) < 2 {
		return substr
	}
	charCount, _ := strconv.Atoi(substr[1:]) // we checked that the second char is digit before
	return strings.Repeat(substr[0:1], charCount)
}
