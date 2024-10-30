package stringBuilder

import (
	"bytes"
	"fmt"
	"testing"
)

type getPadLengthTest struct {
	name          string
	inputString   string
	desiredLength int
	expected      int
}

type getLeftRightPadLengthsTest struct {
	name     string
	numToPad int
	leftPad  int
	rightPad int
}

type getNumCharsToPadTest struct {
	name      string
	padLength int
	numSpaces int
	expected  int
}

type generatePaddedStringTest struct {
	name          string
	inputString   string
	expected      string
	char          rune
	numSpaces     int
	desiredLength int
}

type generatePaddedStringInputs struct {
	inputString   string
	char          rune
	numSpaces     int
	desiredLength int
}

// ##### TEST FUNCTIONS ##### //
func TestGetPadLength(t *testing.T) {
	tests := []getPadLengthTest{
		{
			name:          "basic test no need pad",
			inputString:   "test this",
			desiredLength: 9,
			expected:      0,
		},
		{
			name:          "rune test no need pad",
			inputString:   "👀",
			desiredLength: 1,
			expected:      0,
		},
		{
			name:          "basic test need bad",
			inputString:   "test this",
			desiredLength: 10,
			expected:      1,
		},
		{
			name:          "rune test need pad",
			inputString:   "👀",
			desiredLength: 3,
			expected:      2,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := GetPadLength(test.inputString, test.desiredLength)
			assertNoError(err, t)
			expected := test.expected
			assertValues(got, expected, t)
		})
	}

	negativeTest := getPadLengthTest{
		name:          "test negative",
		inputString:   "👀",
		desiredLength: 0,
		expected:      0,
	}
	t.Run(negativeTest.name, func(t *testing.T) {
		_, err := GetPadLength(negativeTest.inputString, negativeTest.desiredLength)
		assertError(err, t)
	})

}

func TestGetLeftRightPadLengths(t *testing.T) {
	positiveTests := []getLeftRightPadLengthsTest{
		{
			name:     "testing basic even",
			numToPad: 4,
			leftPad:  2,
			rightPad: 2,
		},
		{
			name:     "testing basic odd",
			numToPad: 3,
			leftPad:  1,
			rightPad: 2,
		},
		{
			name:     "testing 0",
			numToPad: 0,
			leftPad:  0,
			rightPad: 0,
		},
	}
	for _, test := range positiveTests {
		t.Run(test.name, func(t *testing.T) {
			gotLeft, gotRight, err := GetLeftRightPadLengths(test.numToPad)
			assertNoError(err, t)
			assertLeftRightPads(gotLeft, gotRight, test.leftPad, test.rightPad, t)
		})
	}

	negativeTest := getLeftRightPadLengthsTest{
		name:     "testing negative",
		numToPad: -1000,
		leftPad:  0,
		rightPad: 0,
	}

	t.Run(negativeTest.name, func(t *testing.T) {
		_, _, err := GetLeftRightPadLengths(negativeTest.numToPad)
		assertError(err, t)
	})

}

func TestGeneratePaddedString(t *testing.T) {
	positiveTests := []generatePaddedStringTest{
		{
			name:          "basic generate padded string test",
			inputString:   "test",
			expected:      "#  test  #",
			char:          '#',
			numSpaces:     2,
			desiredLength: 10,
		},
		{
			name:          "testing emojis",
			inputString:   "test",
			expected:      "🤑  test  🤑",
			char:          '🤑',
			numSpaces:     2,
			desiredLength: 10,
		},
		{
			name:          "testing multiple chars",
			inputString:   "test",
			expected:      "🤑🤑🤑🤑  test  🤑🤑🤑🤑",
			char:          '🤑',
			numSpaces:     2,
			desiredLength: 16,
		},
		{
			name:          "testing odd spaces",
			inputString:   "test",
			expected:      "🤑🤑🤑  test  🤑🤑🤑🤑",
			char:          '🤑',
			numSpaces:     2,
			desiredLength: 15,
		},
		{
			name:          "testing blank string",
			inputString:   "",
			expected:      "🤑🤑🤑🤑🤑🤑    🤑🤑🤑🤑🤑🤑",
			char:          '🤑',
			numSpaces:     2,
			desiredLength: 16,
		},
	}
	for _, test := range positiveTests {
		t.Run(test.name, func(t *testing.T) {
			got, err := GeneratePaddedString(test.inputString,
				test.char, test.numSpaces, test.desiredLength)
			assertStrings(got, test.expected, t)
			assertNoError(err, t)
			assertValues(len(got), len(test.expected), t)
		})
	}
}
func TestGetNumCharsToPad(t *testing.T) {
	positiveTest := getNumCharsToPadTest{
		name:      "postive getNumCharsToPad test",
		padLength: 5,
		numSpaces: 4,
		expected:  1,
	}
	t.Run(positiveTest.name, func(t *testing.T) {
		got, err := getNumCharsToPad(positiveTest.padLength, positiveTest.numSpaces)
		assertNoError(err, t)
		assertValues(got, positiveTest.expected, t)
	})
	negativeTest := getNumCharsToPadTest{
		name:      "negative getNumCharsToPad test",
		padLength: 5,
		numSpaces: 6,
		expected:  0,
	}
	t.Run(negativeTest.name, func(t *testing.T) {
		_, err := getNumCharsToPad(negativeTest.padLength, negativeTest.numSpaces)
		assertError(err, t)
	})
}

func TestPrintString(t *testing.T) {
	t.Run("testing one write to writer", func(t *testing.T) {
		buf := &bytes.Buffer{}
		input := "this is a test string"
		err := PrintString(buf, input)
		assertNoError(err, t)
		got := buf.String()
		expected := input
		assertStrings(got, expected, t)
	})
	t.Run("testing multiple writes to writer", func(t *testing.T) {
		padStringInputs := []generatePaddedStringInputs{
			{
				inputString:   "test",
				char:          '#',
				numSpaces:     2,
				desiredLength: 10,
			},
			{
				inputString:   "test",
				char:          '🤑',
				numSpaces:     4,
				desiredLength: 20,
			},
		}
		buf := &bytes.Buffer{}
		got := ""
		for _, input := range padStringInputs {
			returnedString, err := GeneratePaddedString(input.inputString, input.char, input.numSpaces, input.desiredLength)
			assertNoError(err, t)
			stringToWrite := returnedString + "\n"
			got += stringToWrite
			fmt.Fprint(buf, stringToWrite)
		}
		assertStrings(got, buf.String(), t)

	})

}

// ##### END TEST FUNCTIONS ##### //

// ##### ASSERTIONS #####

func assertStrings(got, expected string, t testing.TB) {
	t.Helper()
	if got != expected {
		t.Errorf("\ngot: %s\nwant: %s", got, expected)
	}
}

func assertValues(got, expected int, t testing.TB) {
	t.Helper()
	if got != expected {
		t.Errorf("\ngot: %d\nwant: %d", got, expected)
	}
}
func assertLeftRightPads(gotLeft, gotRight, wantLeft, wantRight int, t testing.TB) {
	if gotLeft != wantLeft || gotRight != wantRight {
		t.Errorf("\ngotLeft: %d\nwantLeft: %d\ngotRight %d\nwantRight %d\n",
			gotLeft, gotRight, wantLeft, wantRight)
	}
}

func assertNoError(err error, t testing.TB) {
	if err != nil {
		t.Errorf("Unexpected error, %s", err.Error())
	}
}

func assertError(err error, t testing.TB) {
	if err == nil {
		t.Errorf("expected error, didn't get one")
	}
}

// ##### END ASSERTIONS #####
