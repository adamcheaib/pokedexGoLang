package main

import (
	"strings"
)

func CleanInput(word string) []string {
	splitWords := strings.Split(word, " ")

	if len(splitWords) >= 2 {
		return []string{
			strings.ToLower(splitWords[0]),
			strings.ToLower(splitWords[1])}
	}
	return []string{strings.ToLower(splitWords[0])}
}
