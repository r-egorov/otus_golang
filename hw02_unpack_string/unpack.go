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

	runes := []rune(inputStr)
	for pos := 0; pos < len(runes); {
		substr, err := findNextSubstr(runes[pos:])
		if err != nil {
			return "", err
		}
		b.WriteString(unpackSubstr(substr))
		pos += len(substr)
	}
	return b.String(), nil
}

func findNextSubstr(runes []rune) ([]rune, error) {
	if unicode.IsDigit(runes[0]) {
		return nil, ErrInvalidString
	}

	screenOffset := 0
	if runes[0] == '\\' {
		canBeScreened := len(runes) > 1 && (unicode.IsDigit(runes[1]) || runes[1] == '\\')
		if !canBeScreened {
			return nil, ErrInvalidString
		}
		screenOffset = 1
	}

	isSubstrOneChar := len(runes) < 2+screenOffset || !unicode.IsDigit(runes[1+screenOffset])
	if isSubstrOneChar {
		return runes[:1+screenOffset], nil
	}

	isMultipleCyphersInSubstr := len(runes) > 2+screenOffset && unicode.IsDigit(runes[2+screenOffset])
	if isMultipleCyphersInSubstr {
		return nil, ErrInvalidString
	}
	return runes[:2+screenOffset], nil
}

func unpackSubstr(runes []rune) string {
	if runes[0] == '\\' {
		runes = runes[1:]
	}
	if len(runes) < 2 {
		return string(runes)
	}
	charCount, _ := strconv.Atoi(string(runes[1])) // we checked that the second char is digit before
	return strings.Repeat(string(runes[0:1]), charCount)
}
