package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
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
	rawTag, needToBeValidated := field.Tag.Lookup("validate")
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
		err = processSlice(field, tags)
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
	tagSplitted := strings.SplitN(rawTags, ":", 2)
	if len(tagSplitted) != 2 {
		return nil, ErrTagInvalid
	}
	tags[tagSplitted[0]] = tagSplitted[1]
	return tags, nil
}

func processSlice(sl reflect.StructField, tags map[string]string) error {
	elemType := sl.Type.Elem()
	switch elemType.Kind() { //nolint:exhaustive
	case reflect.String:
		fmt.Println("slice of strings")
	case reflect.Int:
		fmt.Println("slice of ints")
	default:
		fmt.Println("slice not supported")
	}
	return nil
}
