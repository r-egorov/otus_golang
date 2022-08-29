package hw09structvalidator

import (
	"strconv"
	"strings"
)

func processInt(fieldValue int, tags map[string]string) error {
	var err error
	for tagKey := range tags {
		tagValue := tags[tagKey]
		switch tagKey {
		case minTag:
			err = checkIntMinTag(fieldValue, tagValue)
		case maxTag:
			err = checkIntMaxTag(fieldValue, tagValue)
		case inTag:
			err = checkIntInSet(fieldValue, strings.Split(tagValue, ","))
		default:
			err = ErrTagInvalid
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func checkIntMinTag(fieldValue int, tagValue string) error {
	minimum, atoiErr := strconv.Atoi(tagValue)
	if atoiErr != nil {
		return ErrTagInvalid
	}
	if fieldValue < minimum {
		return ErrIntUnderMinimum
	}
	return nil
}

func checkIntMaxTag(fieldValue int, tagValue string) error {
	maximum, atoiErr := strconv.Atoi(tagValue)
	if atoiErr != nil {
		return ErrTagInvalid
	}
	if fieldValue > maximum {
		return ErrIntOverMaximum
	}
	return nil
}

func checkIntInSet(fieldValue int, intStrs []string) error {
	if len(intStrs) < 1 {
		return ErrTagInvalid
	}
	for i := 0; i < len(intStrs); i++ {
		atoied, atoiErr := strconv.Atoi(intStrs[i])
		if atoiErr != nil {
			return ErrTagInvalid
		}
		if fieldValue == atoied {
			return nil
		}
	}
	return ErrIntNotInSet
}
