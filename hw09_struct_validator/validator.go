package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var (
	ErrNotAStruct      = errors.New("parameter is not a struct")
	ErrTagInvalid      = errors.New("invalid tag")
	ErrTypeUnsupported = errors.New("unsupported type")

	ErrStringLenInvalid    = errors.New("invalid string length")
	ErrStringNotInSet      = errors.New("string not in set")
	ErrStringNotRegexplike = errors.New("string not regexplike")

	ErrIntUnderMinimum = errors.New("int less than minimum")
	ErrIntOverMaximum  = errors.New("int greater than maximuim")
	ErrIntNotInSet     = errors.New("int not in set")
)

type ValidationError struct {
	Field string
	Err   error
}

func (ve ValidationError) Error() string {
	return fmt.Sprintf("{%s: %s}", ve.Field, ve.Err)
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	if len(v) < 1 {
		return "no validation errors"
	}
	builder := strings.Builder{}
	builder.WriteString("Validation errors: ")
	for i, err := range v {
		builder.WriteString(err.Error())
		if i < len(v)-1 {
			builder.WriteString(", ")
		}
	}
	return builder.String()
}

func Validate(v interface{}) error {
	valOfStruct := reflect.ValueOf(v)
	typeOfStruct := reflect.TypeOf(v)
	errs := ValidationErrors{}

	if valOfStruct.Kind() != reflect.Struct {
		return ErrNotAStruct
	}

	for i := 0; i < typeOfStruct.NumField(); i++ {
		err := validateField(typeOfStruct.Field(i), valOfStruct.Field(i))
		if err != nil {
			errs = append(errs, *err)
		}
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}

func validateField(field reflect.StructField, fieldValue reflect.Value) *ValidationError {
	tag, needToBeValidated := field.Tag.Lookup("validate")
	if !needToBeValidated {
		return nil
	}

	var err error
	switch field.Type.Kind() {
	case reflect.Slice:
		err = processSlice(field, tag)
	case reflect.String:
		err = processString(fieldValue.String(), tag)
	case reflect.Int:
		err = processInt(int(fieldValue.Int()), tag)
	default:
		err = ErrTypeUnsupported
	}

	if err != nil {
		validationErr := &ValidationError{
			Field: field.Name,
			Err:   err,
		}
		return validationErr
	}
	return nil
}

func processSlice(sl reflect.StructField, tag string) error {
	elemType := sl.Type.Elem()
	switch elemType.Kind() {
	case reflect.String:
		fmt.Println("slice of strings")
	case reflect.Int:
		fmt.Println("slice of ints")
	default:
		fmt.Println("slice not supported")
	}
	return nil
}

func processString(str string, tag string) error {
	tagSplitted := strings.SplitN(tag, ":", 2)
	if len(tagSplitted) != 2 {
		return ErrTagInvalid
	}
	tagKey := tagSplitted[0]
	tagValue := tagSplitted[1]

	var err error
	switch tagKey {
	case "len":
		expectedLen, atoiErr := strconv.Atoi(tagValue)
		if atoiErr != nil {
			return ErrTagInvalid
		}
		err = checkStringLen(str, expectedLen)
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
	return nil
}

func checkStringLen(str string, expectedLen int) error {
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

func processInt(fieldValue int, tag string) error {
	tagSplitted := strings.SplitN(tag, ":", 2)
	if len(tagSplitted) != 2 {
		return ErrTagInvalid
	}
	tagKey := tagSplitted[0]
	tagValue := tagSplitted[1]

	var err error
	switch tagKey {
	case "min":
		minimum, atoiErr := strconv.Atoi(tagValue)
		if atoiErr != nil {
			return ErrTagInvalid
		}
		if fieldValue < minimum {
			return ErrIntUnderMinimum
		}
	case "max":
		maximum, atoiErr := strconv.Atoi(tagValue)
		if atoiErr != nil {
			return ErrTagInvalid
		}
		if fieldValue > maximum {
			return ErrIntOverMaximum
		}
	case "in":
		intsStrs := strings.Split(tagValue, ",")
		if len(intsStrs) < 1 {
			return ErrTagInvalid
		}
		for i := 0; i < len(intsStrs); i++ {
			atoied, atoiErr := strconv.Atoi(intsStrs[i])
			if atoiErr != nil {
				return ErrTagInvalid
			}
			if fieldValue == atoied {
				return nil
			}
		}
		return ErrIntNotInSet
	default:
		err = ErrTagInvalid
	}

	if err != nil {
		return err
	}
	return nil
}
