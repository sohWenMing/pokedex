package utils

import (
	"fmt"
	"strings"
)

func CleanInput(text string) []string {

	splitWords := strings.Fields(text)
	returned := make([]string, 0, len(splitWords))
	for _, word := range splitWords {
		returned = append(returned, strings.ToLower(strings.TrimSpace(word)))
	}
	return returned
}

func CleanLineAndAddNewLine(text string) string {
	cleanedString := strings.TrimSuffix(strings.TrimPrefix(text, "\n"), "\n")
	return fmt.Sprintf("%s\n", cleanedString)
}
