package main

import (
	"bytes"
	"reflect"
	"testing"
)

func TestOk(t *testing.T) {
	var cases = []struct {
		expected int
		input    string
	}{
		{
			expected: 3,
			input: "1 2 + =",
		},

		{
			expected: 0,
			input: "1 0 * =",
		},

		{
			expected: 1,
			input: "3 2 - =",
		},

		{
			expected: 2,
			input: "6 3 / =",
		},

		{
			expected: 15,
			input: "1 2 3 4 + * + =",
		},

		{
			expected: 21,
			input: " 1 2 + 3 4 + * =",
		},
	}
	for _, item := range cases {
		result, _ := calc(bytes.NewBufferString(item.input))
		if !reflect.DeepEqual(result, item.expected) {
			t.Error("expected", item.expected, "have", result)
		}
	}
}

func TestFail(t *testing.T) {
	var cases = []string{"a + 1 =", "1 2 + *"}
	for _, item := range cases {
		_, err := calc(bytes.NewBufferString(item))
		if err == nil {
			t.Errorf("Test FAIL failed: expected error")
		}
	}
}
