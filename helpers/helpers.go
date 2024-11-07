package helpers

import "strings"

func ToLowerAndTrim(input string) (output string) {
	output = strings.ToLower(strings.TrimSpace(input))
	return output
}
