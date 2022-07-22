package hw09structvalidator

import (
	"regexp"
	"strconv"
	"strings"
)

func processString(str string, tags map[string]string) error {
	var err error
	for tagKey := range tags {
		tagValue := tags[tagKey]
		switch tagKey {
		case "len":
			err = checkStringLenTag(str, tagValue)
		case "in":
			err = checkStringInSet(str, strings.Split(tagValue, ","))
		case "regexp":
			err = checkStringRegexplike(str, tagValue)
		default:
			err = ErrTagInvalid
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func checkStringLenTag(str, tagValue string) error {
	expectedLen, atoiErr := strconv.Atoi(tagValue)
	if atoiErr != nil {
		return ErrTagInvalid
	}
	if len([]rune(str)) != expectedLen {
		return ErrStringLenInvalid
	}
	return nil
}

func checkStringInSet(str string, set []string) error {
	if len(set) < 1 {
		return ErrTagInvalid
	}
	for _, s := range set {
		if str == s {
			return nil
		}
	}
	return ErrStringNotInSet
}

func checkStringRegexplike(str string, regexStr string) error {
	regex, err := regexp.Compile(regexStr)
	if err != nil {
		return ErrTagInvalid
	}

	matched := regex.MatchString(str)
	if !matched {
		return ErrStringNotRegexplike
	}
	return nil
}
