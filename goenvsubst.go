package goenvsubst

import (
	"os"
	"reflect"
	"strings"
)

// Do recursively walks through any Go data structure (structs, slices, maps, arrays, pointers)
// and replaces environment variable references in string values with their actual values
// from the environment. Environment variables should be in the format $VAR_NAME.
// Supports top-level and nested: structs, slices, arrays, maps, and pointers.
func Do(v any) error {
	return doValue(reflect.ValueOf(v))
}

// doValue recursively processes reflect.Value to expand environment variables
func doValue(v reflect.Value) error {
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.String:
		return doString(v)
	case reflect.Struct:
		return doStruct(v)
	case reflect.Slice, reflect.Array:
		return doSliceArray(v)
	case reflect.Map:
		return doMap(v)
	}

	return nil
}

// doString processes string values for environment variable expansion
func doString(v reflect.Value) error {
	if v.CanSet() {
		v.SetString(expandEnvVar(v.String()))
	}
	return nil
}

// doStruct processes struct values recursively
func doStruct(v reflect.Value) error {
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.CanSet() {
			if err := doValue(field); err != nil {
				return err
			}
		}
	}
	return nil
}

// doSliceArray processes slice and array values recursively
func doSliceArray(v reflect.Value) error {
	for i := 0; i < v.Len(); i++ {
		if err := doValue(v.Index(i)); err != nil {
			return err
		}
	}
	return nil
}

// doMap processes map values recursively
func doMap(v reflect.Value) error {
	for _, key := range v.MapKeys() {
		mapValue := v.MapIndex(key)
		// For maps, we need to create a new value, modify it, and set it back
		if mapValue.Kind() == reflect.String {
			original := mapValue.String()
			expanded := expandEnvVar(original)
			if expanded != original {
				v.SetMapIndex(key, reflect.ValueOf(expanded))
			}
		} else {
			// For non-string values, create a copy and recurse
			newValue := reflect.New(mapValue.Type()).Elem()
			newValue.Set(mapValue)
			if err := doValue(newValue); err != nil {
				return err
			}
			v.SetMapIndex(key, newValue)
		}
	}
	return nil
}

// expandEnvVar replaces environment variable references in the format $VAR_NAME
// with their actual values from the environment. Returns empty string for
// missing or empty environment variables.
func expandEnvVar(s string) string {
	if !strings.HasPrefix(s, "$") {
		return s
	}

	// Remove the $ prefix to get the variable name
	varName := strings.TrimPrefix(s, "$")

	// Get the environment variable value
	return os.Getenv(varName)
}
