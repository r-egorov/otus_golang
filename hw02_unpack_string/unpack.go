package hw02unpackstring

import (
	"errors"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(inputStr string) (string, error) {
	//var b strings.Builder
	//var char string
	//var charCount int
	//var err error
	//
	//char = ""
	//charCount = 1
	//for pos, r := range inputStr {
	//	if char == "" {
	//		if unicode.IsDigit(r) {
	//			return "", ErrInvalidString
	//		}
	//		char = inputStr[pos : pos+1]
	//	} else {
	//		if unicode.IsDigit(r) {
	//			charCount, err = strconv.Atoi(inputStr[pos : pos+1])
	//			if err != nil {
	//				return "", ErrInvalidString
	//			}
	//			b.WriteString(strings.Repeat(char, charCount))
	//			char = ""
	//		} else {
	//			b.WriteString(char)
	//			char = inputStr[pos : pos+1]
	//		}
	//	}
	//}
	//if char != "" {
	//	b.WriteString(char)
	//}
	//return b.String(), nil

	var b strings.Builder

	for pos := 0; pos < len(inputStr); {
		substr, err := findSubstr(inputStr[pos:])
		if err != nil {
			return "", err
		}
		unpackedSubstr := unpackSubstr(substr)
		b.WriteString(unpackedSubstr)
		pos += len(substr)
	}
	return b.String(), nil
}

func findSubstr(inputStr string) (string, error) {
	runes := []rune(inputStr)
	isSubstrOneChar := len(runes) < 2 || !unicode.IsDigit(runes[1])
	if isSubstrOneChar {
		return inputStr[:1], nil
	} else {
		return inputStr[:2], nil
	}
}

func unpackSubstr(substr string) string {
	return ""
}
