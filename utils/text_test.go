package utils

import (
	"testing"
)

func TestGenerateID(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Hello World", "hello_world"},
		{"GoLang 123", "golang_123"},
		{"Test-Case_Example", "test_case_example"},
		{"Special@#%Characters", "specialcharacters"},
		{"   Leading and trailing spaces   ", "leading_and_trailing_spaces"},
		{"MixedCASE Input", "mixedcase_input"},
	}

	for _, test := range tests {
		result := GenerateID(test.input)
		if result != test.expected {
			t.Errorf("GenerateID(%q) = %q; expected %q", test.input, result, test.expected)
		}
	}
}
