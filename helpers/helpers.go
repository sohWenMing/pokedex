package helpers

import "strings"

func ToLowerAndTrim(input string) (output string) {
	output = strings.ToLower(strings.TrimSpace(input))
	return output
}

func ReplaceSpaceWithChar(inputString, replacementChar string) (outputString string) {
	outputString = StringReplace(inputString, " ", replacementChar)
	return outputString
}

func StringReplace(inputString, stringToReplace, replaceWithString string) (outputString string) {
	outputString = strings.Replace(inputString, stringToReplace, replaceWithString, -1)
	return outputString
}
