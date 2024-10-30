package stringBuilder

import "unicode/utf8"

func GetPadLength(inputString string, desiredLength int) (numToPad int) {
	numChars := utf8.RuneCountInString(inputString)
	if desiredLength-numChars < 0 {
		return 0
	}
	return desiredLength - numChars
}

func GetLeftRightPadLengths(numToPad int) (leftPad, rightPad int) {
	if numToPad <= 0 {
		return
	}
	switch numToPad%2 == 0 {
	case true:
		leftPad = numToPad / 2
		rightPad = numToPad / 2
	case false:
		leftPad = numToPad / 2
		rightPad = (numToPad / 2) + 1
	}
	return leftPad, rightPad
}
