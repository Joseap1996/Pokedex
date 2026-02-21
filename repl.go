package main

import (
	"strings"
)

func cleanInput(text string) []string {
	input := strings.ToLower(text)
	words := strings.Fields(input)

	return words
}
