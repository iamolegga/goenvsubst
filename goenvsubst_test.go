package goenvsubst_test

import (
	"os"
	"reflect"
	"testing"

	"github.com/iamolegga/goenvsubst"
)

var tests = []struct {
	name     string
	input    any
	expected any
}{
	{
		name:     "simple string substitution",
		input:    &struct{ Value string }{"$TEST_VAR"},
		expected: &struct{ Value string }{"test_value"},
	},
	{
		name:     "string without substitution",
		input:    &struct{ Value string }{"no_substitution"},
		expected: &struct{ Value string }{"no_substitution"},
	},
	{
		name:     "missing environment variable",
		input:    &struct{ Value string }{"$MISSING_VAR"},
		expected: &struct{ Value string }{""},
	},
	{
		name:     "empty environment variable",
		input:    &struct{ Value string }{"$EMPTY_VAR"},
		expected: &struct{ Value string }{""},
	},
	{
		name: "multiple fields",
		input: &struct {
			A, B, C string
		}{"$TEST_VAR", "static", "$ANOTHER_VAR"},
		expected: &struct {
			A, B, C string
		}{"test_value", "static", "another_value"},
	},
	{
		name: "nested struct",
		input: &struct {
			Outer string
			Inner struct{ Value string }
		}{"$TEST_VAR", struct{ Value string }{"$ANOTHER_VAR"}},
		expected: &struct {
			Outer string
			Inner struct{ Value string }
		}{"test_value", struct{ Value string }{"another_value"}},
	},
	{
		name:     "slice of strings",
		input:    &struct{ Values []string }{[]string{"$TEST_VAR", "static", "$ANOTHER_VAR"}},
		expected: &struct{ Values []string }{[]string{"test_value", "static", "another_value"}},
	},
	{
		name: "slice of structs",
		input: &struct{ Items []struct{ Name string } }{
			[]struct{ Name string }{{"$TEST_VAR"}, {"static"}, {"$ANOTHER_VAR"}},
		},
		expected: &struct{ Items []struct{ Name string } }{
			[]struct{ Name string }{{"test_value"}, {"static"}, {"another_value"}},
		},
	},
	{
		name: "map with string values",
		input: &struct{ Data map[string]string }{
			map[string]string{"key1": "$TEST_VAR", "key2": "static", "key3": "$ANOTHER_VAR"},
		},
		expected: &struct{ Data map[string]string }{
			map[string]string{"key1": "test_value", "key2": "static", "key3": "another_value"},
		},
	},
	{
		name: "map keys are not replaced",
		input: &struct{ Data map[string]string }{
			map[string]string{"$TEST_VAR": "value1", "static_key": "$ANOTHER_VAR", "$MISSING_VAR": "value3"},
		},
		expected: &struct{ Data map[string]string }{
			map[string]string{"$TEST_VAR": "value1", "static_key": "another_value", "$MISSING_VAR": "value3"},
		},
	},
	{
		name: "mixed types (int, bool, string)",
		input: &struct {
			S string
			I int
			B bool
		}{"$TEST_VAR", 42, true},
		expected: &struct {
			S string
			I int
			B bool
		}{"test_value", 42, true},
	},
	{
		name:     "pointer to struct",
		input:    &struct{ Ptr *struct{ Value string } }{&struct{ Value string }{"$TEST_VAR"}},
		expected: &struct{ Ptr *struct{ Value string } }{&struct{ Value string }{"test_value"}},
	},
	{
		name:     "nil pointer",
		input:    &struct{ Ptr *struct{ Value string } }{nil},
		expected: &struct{ Ptr *struct{ Value string } }{nil},
	},
	// Top-level slice tests
	{
		name:     "top-level slice of strings",
		input:    &[]string{"$TEST_VAR", "static", "$ANOTHER_VAR", "$MISSING_VAR"},
		expected: &[]string{"test_value", "static", "another_value", ""},
	},
	{
		name:     "top-level empty slice",
		input:    &[]string{},
		expected: &[]string{},
	},
	{
		name:     "top-level slice with empty strings",
		input:    &[]string{"", "$TEST_VAR", ""},
		expected: &[]string{"", "test_value", ""},
	},
	{
		name:     "top-level slice of structs",
		input:    &[]struct{ Name string }{{"$TEST_VAR"}, {"static"}, {"$ANOTHER_VAR"}},
		expected: &[]struct{ Name string }{{"test_value"}, {"static"}, {"another_value"}},
	},
	{
		name:     "top-level nested slice",
		input:    &[][]string{{"$TEST_VAR", "static"}, {"$ANOTHER_VAR", "more"}},
		expected: &[][]string{{"test_value", "static"}, {"another_value", "more"}},
	},
	// Top-level map tests
	{
		name:     "top-level map string to string",
		input:    &map[string]string{"key1": "$TEST_VAR", "key2": "static", "key3": "$ANOTHER_VAR"},
		expected: &map[string]string{"key1": "test_value", "key2": "static", "key3": "another_value"},
	},
	{
		name:     "top-level empty map",
		input:    &map[string]string{},
		expected: &map[string]string{},
	},
	{
		name:     "top-level map keys with env vars (should not be replaced)",
		input:    &map[string]string{"$TEST_VAR": "value1", "static_key": "$ANOTHER_VAR"},
		expected: &map[string]string{"$TEST_VAR": "value1", "static_key": "another_value"},
	},
	{
		name:     "top-level map with struct values",
		input:    &map[string]struct{ Value string }{"key1": {"$TEST_VAR"}, "key2": {"$ANOTHER_VAR"}},
		expected: &map[string]struct{ Value string }{"key1": {"test_value"}, "key2": {"another_value"}},
	},
	{
		name:     "top-level map with slice values",
		input:    &map[string][]string{"key1": {"$TEST_VAR", "static"}, "key2": {"$ANOTHER_VAR"}},
		expected: &map[string][]string{"key1": {"test_value", "static"}, "key2": {"another_value"}},
	},
	{
		name:     "top-level nested map",
		input:    &map[string]map[string]string{"outer": {"inner": "$TEST_VAR"}},
		expected: &map[string]map[string]string{"outer": {"inner": "test_value"}},
	},
	// Top-level array tests
	{
		name:     "top-level array of strings",
		input:    &[3]string{"$TEST_VAR", "static", "$ANOTHER_VAR"},
		expected: &[3]string{"test_value", "static", "another_value"},
	},
	{
		name:     "top-level array of structs",
		input:    &[2]struct{ Name string }{{"$TEST_VAR"}, {"$ANOTHER_VAR"}},
		expected: &[2]struct{ Name string }{{"test_value"}, {"another_value"}},
	},
	// Top-level pointer tests
	{
		name: "top-level pointer to string",
		input: func() *string {
			s := "$TEST_VAR"
			return &s
		}(),
		expected: func() *string {
			s := "test_value"
			return &s
		}(),
	},
	{
		name: "top-level pointer to struct",
		input: func() *struct{ Value string } {
			return &struct{ Value string }{"$TEST_VAR"}
		}(),
		expected: func() *struct{ Value string } {
			return &struct{ Value string }{"test_value"}
		}(),
	},
	{
		name: "top-level nil pointer",
		input: func() *string {
			return nil
		}(),
		expected: func() *string {
			return nil
		}(),
	},

	// Complex nested structures
	{
		name: "complex nested structure",
		input: &map[string][]struct {
			Name   string
			Values []string
			Meta   map[string]string
		}{
			"group1": {
				{
					Name:   "$TEST_VAR",
					Values: []string{"$ANOTHER_VAR", "static"},
					Meta:   map[string]string{"key": "$TEST_VAR"},
				},
			},
		},
		expected: &map[string][]struct {
			Name   string
			Values []string
			Meta   map[string]string
		}{
			"group1": {
				{
					Name:   "test_value",
					Values: []string{"another_value", "static"},
					Meta:   map[string]string{"key": "test_value"},
				},
			},
		},
	},
	// Additional edge cases
	{
		name:     "slice of pointers",
		input:    &[]*string{func() *string { s := "$TEST_VAR"; return &s }(), func() *string { s := "$ANOTHER_VAR"; return &s }(), nil},
		expected: &[]*string{func() *string { s := "test_value"; return &s }(), func() *string { s := "another_value"; return &s }(), nil},
	},
	{
		name:     "map with pointer values",
		input:    &map[string]*string{"key1": func() *string { s := "$TEST_VAR"; return &s }(), "key2": nil},
		expected: &map[string]*string{"key1": func() *string { s := "test_value"; return &s }(), "key2": nil},
	},
}

func TestDo(t *testing.T) {
	// Set up test environment variables
	os.Setenv("TEST_VAR", "test_value")
	os.Setenv("ANOTHER_VAR", "another_value")
	os.Setenv("EMPTY_VAR", "")
	defer func() {
		os.Unsetenv("TEST_VAR")
		os.Unsetenv("ANOTHER_VAR")
		os.Unsetenv("EMPTY_VAR")
	}()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := goenvsubst.Do(tt.input)
			if err != nil {
				t.Errorf("Do() error = %v", err)
				return
			}

			// Compare the results
			if !deepEqual(tt.input, tt.expected) {
				t.Errorf("Do() result mismatch.\nGot: %+v\nWant: %+v", tt.input, tt.expected)
			}
		})
	}
}

// deepEqual performs a deep comparison of two interfaces
// This is a simplified version for our test cases
func deepEqual(a, b any) bool {
	return compareValues(reflect.ValueOf(a), reflect.ValueOf(b))
}

func compareValues(a, b reflect.Value) bool {
	if a.Type() != b.Type() {
		return false
	}

	switch a.Kind() {
	case reflect.Ptr:
		return comparePointers(a, b)
	case reflect.Struct:
		return compareStructs(a, b)
	case reflect.Slice:
		return compareSlices(a, b)
	case reflect.Map:
		return compareMaps(a, b)
	default:
		return a.Interface() == b.Interface()
	}
}

// comparePointers compares two pointer values
func comparePointers(a, b reflect.Value) bool {
	if a.IsNil() && b.IsNil() {
		return true
	}
	if a.IsNil() || b.IsNil() {
		return false
	}
	return compareValues(a.Elem(), b.Elem())
}

// compareStructs compares two struct values field by field
func compareStructs(a, b reflect.Value) bool {
	for i := 0; i < a.NumField(); i++ {
		if !compareValues(a.Field(i), b.Field(i)) {
			return false
		}
	}
	return true
}

// compareSlices compares two slice values element by element
func compareSlices(a, b reflect.Value) bool {
	if a.Len() != b.Len() {
		return false
	}
	for i := 0; i < a.Len(); i++ {
		if !compareValues(a.Index(i), b.Index(i)) {
			return false
		}
	}
	return true
}

// compareMaps compares two map values key by key
func compareMaps(a, b reflect.Value) bool {
	if a.Len() != b.Len() {
		return false
	}
	for _, key := range a.MapKeys() {
		aVal := a.MapIndex(key)
		bVal := b.MapIndex(key)
		if !bVal.IsValid() || !compareValues(aVal, bVal) {
			return false
		}
	}
	return true
}
