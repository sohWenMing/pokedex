package helpers

import (
	"fmt"
	"strings"
)

func PrintPrompt(input string) {
	fmt.Printf("%s> ", input)
}

func TrimToLowerString(input string) string {
	input = strings.TrimSpace(input)
	input = strings.ToLower(input)
	return input
}
