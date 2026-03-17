package main

import (
	"strings"
)

func cleanInput(text string) []string {
	result := []string{}
	result = strings.Fields(text)
	for i, word := range result {
		result[i] = strings.ToLower(word)
	}
	return result
}
