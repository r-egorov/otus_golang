package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

const (
	validateTagKey = "validate"
	lenTag         = "len"
	inTag          = "in"
	regexpTag      = "regexp"
	minTag         = "min"
	maxTag         = "max"
)

var (
	ErrNotAStruct      = errors.New("parameter is not a struct")
	ErrTagInvalid      = errors.New("invalid tag")
	ErrTypeUnsupported = errors.New("unsupported type")

	ErrStringLenInvalid    = errors.New("invalid string length")
	ErrStringNotInSet      = errors.New("string not in set")
	ErrStringNotRegexplike = errors.New("string not regexplike")

	ErrIntUnderMinimum = errors.New("int less than Minimum")
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
		if valOfStruct.Field(i).CanInterface() { // if false, field is unexported
			err := validateField(typeOfStruct.Field(i), valOfStruct.Field(i))
			if err != nil {
				errs = append(errs, *err)
			}
		}
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}

func validateField(field reflect.StructField, fieldValue reflect.Value) *ValidationError {
	rawTag, needToBeValidated := field.Tag.Lookup(validateTagKey)
	if !needToBeValidated {
		return nil
	}

	tags, err := parseTags(rawTag)
	if err != nil {
		return &ValidationError{
			Field: field.Name,
			Err:   err,
		}
	}

	switch field.Type.Kind() { //nolint:exhaustive
	case reflect.Slice:
		err = processSlice(field, fieldValue, tags)
	case reflect.String:
		err = processString(fieldValue.String(), tags)
	case reflect.Int:
		err = processInt(int(fieldValue.Int()), tags)
	default:
		err = ErrTypeUnsupported
	}

	if err != nil {
		return &ValidationError{
			Field: field.Name,
			Err:   err,
		}
	}
	return nil
}

func parseTags(rawTags string) (map[string]string, error) {
	tags := make(map[string]string)
	rawTagsSplitted := strings.Split(rawTags, "|")
	for _, rawTag := range rawTagsSplitted {
		tagSplitted := strings.SplitN(rawTag, ":", 2)
		tags[tagSplitted[0]] = tagSplitted[1]
	}
	return tags, nil
}
