package stringBuilder

import (
	"fmt"
	"io"
	"strings"
	"unicode/utf8"
)

func GetPadLength(inputString string, desiredLength int) (numToPad int, err error) {
	numChars := utf8.RuneCountInString(inputString)
	if desiredLength-numChars < 0 {
		return 0, fmt.Errorf("desiredLength cannot be less than length of input string")
	}
	return desiredLength - numChars, nil
}

func GetLeftRightPadLengths(numToPad int) (leftPad, rightPad int, err error) {
	if numToPad < 0 {
		return 0, 0, fmt.Errorf("numToPad cannot be negative")
	}
	if numToPad == 0 {
		return 0, 0, nil
	}
	switch numToPad%2 == 0 {
	case true:
		leftPad = numToPad / 2
		rightPad = numToPad / 2
	case false:
		leftPad = numToPad / 2
		rightPad = (numToPad / 2) + 1
	}
	return leftPad, rightPad, nil
}

/*
	if i know the left pad, and the right pad
	determine the character at the edges
	determine the number of spaces
*/

func getNumCharsToPad(padLength, numSpaces int) (numCharsToPad int, err error) {
	if numSpaces > padLength {
		return 0, fmt.Errorf("numSpaces cannot be greater than padLength")
	}
	return padLength - numSpaces, nil

}

func GeneratePaddedString(inputString string, char rune, numSpaces, desiredLength int) (output string, err error) {
	totalPadLength, err := GetPadLength(inputString, desiredLength)
	if err != nil {
		return "", err
	}
	leftPadLength, rightPadlength, err := GetLeftRightPadLengths(totalPadLength)
	if err != nil {
		return "", err
	}
	leftNumChars, err := getNumCharsToPad(leftPadLength, numSpaces)
	if err != nil {
		return "", err
	}
	rightNumChars, err := getNumCharsToPad(rightPadlength, numSpaces)
	if err != nil {
		return "", err
	}

	leftPadString := strings.Repeat(string(char), leftNumChars) + strings.Repeat(" ", numSpaces)
	rightPadString := strings.Repeat(" ", numSpaces) + strings.Repeat(string(char), rightNumChars)
	return leftPadString + inputString + rightPadString, nil

}

func PrintString(w io.Writer, input string) (err error) {
	_, printErr := fmt.Fprint(w, input)
	if printErr != nil {
		return printErr
	}
	return nil
}
