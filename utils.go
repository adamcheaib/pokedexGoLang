package main

import (
	"strings"
)

func CleanInput(word string) []string {
	splitWords := strings.Split(word, " ")
	return []string{strings.ToLower(splitWords[0])}
}
