package hw09structvalidator

import "reflect"

func processSlice(
	sl reflect.StructField,
	fieldValue reflect.Value,
	tags map[string]string,
) error {
	var err error
	elemType := sl.Type.Elem()
	switch elemType.Kind() { //nolint:exhaustive
	case reflect.String:
		slice, ok := fieldValue.Interface().([]string)
		if !ok {
			return ErrTypeUnsupported
		}
		err = processStringSlice(slice, tags)
	case reflect.Int:
		slice, ok := fieldValue.Interface().([]int)
		if !ok {
			return ErrTypeUnsupported
		}
		err = processIntSlice(slice, tags)
	default:
		err = ErrTypeUnsupported
	}
	if err != nil {
		return err
	}
	return nil
}

func processStringSlice(slice []string, tags map[string]string) error {
	for _, str := range slice {
		err := processString(str, tags)
		if err != nil {
			return err
		}
	}
	return nil
}

func processIntSlice(slice []int, tags map[string]string) error {
	for _, intElem := range slice {
		err := processInt(intElem, tags)
		if err != nil {
			return err
		}
	}
	return nil
}
