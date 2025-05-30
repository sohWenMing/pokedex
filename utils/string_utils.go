package utils

import "strings"

func CleanInput(text string) []string {

	splitWords := strings.Fields(text)
	returned := make([]string, 0, len(splitWords))
	for _, word := range splitWords {
		returned = append(returned, strings.ToLower(strings.TrimSpace(word)))
	}
	return returned
}
