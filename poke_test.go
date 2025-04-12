package main

import (
	"fmt"
	"testing"
)

func TestCleanInput(test *testing.T) {
	cases := []struct {
		input    string
		expected []string
		length   int
	}{
		{
			input:    "hello world rastaman",
			expected: []string{"hello", "world"},
			length:   2,
		},
		{
			input:    "Hi",
			expected: []string{"hi"},
			length:   1,
		},
	}

	for _, currentCase := range cases {
		actual := CleanInput(currentCase.input)
		if len(actual) != currentCase.length {
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
