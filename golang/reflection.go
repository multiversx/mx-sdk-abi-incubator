package abi

import "reflect"

func deepClone(src interface{}) interface{} {
	srcValue := reflect.ValueOf(src)
	if srcValue.Kind() != reflect.Struct {
		// If it's not a struct, return as is
		return src
	}

	// Create a new instance of the same type as the source
	dest := reflect.New(srcValue.Type()).Elem()

	for i := 0; i < srcValue.NumField(); i++ {
		srcField := srcValue.Field(i)
		destField := dest.Field(i)

		if srcField.Kind() == reflect.Struct {
			// Recursively clone nested structs
			destField.Set(reflect.ValueOf(deepClone(srcField.Interface())))
		} else {
			// Copy non-struct fields
			destField.Set(srcField)
		}
	}

	return dest.Interface()
}
