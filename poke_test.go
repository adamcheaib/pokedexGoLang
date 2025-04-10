package main

import (
	"fmt"
	"testing"
)

func TestCleanInput(test *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "hello world",
			expected: []string{"hello"},
		},
	}

	for _, currentCase := range cases {
		actual := CleanInput(currentCase.input)
		if len(actual) != len(currentCase.expected) {
			test.Errorf("Expected %v, got %v", len(currentCase.expected), len(actual))
			return
		}

		for i := range currentCase.expected {
			word := actual[i]
			expectedWord := currentCase.expected[i]
			fmt.Println(word, expectedWord)
			if string(word) != expectedWord {
				test.Fatalf("Expected %v, got %v", expectedWord, word)
			}
		}
	}
}
