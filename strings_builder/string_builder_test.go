package stringBuilder

import (
	"bytes"
	"testing"
)

type GetPadLengthTest struct {
	name          string
	inputString   string
	desiredLength int
	expected      int
}

type GetLeftRightPadLengthsTest struct {
	name     string
	numToPad int
	leftPad  int
	rightPad int
}

type GetNumCharsToPadTest struct {
	name      string
	padLength int
	numSpaces int
	expected  int
}

type GeneratePaddedStringTest struct {
	name          string
	inputString   string
	expected      string
	char          rune
	numSpaces     int
	desiredLength int
}

// ##### TEST FUNCTIONS ##### //
func TestGetPadLength(t *testing.T) {
	tests := []GetPadLengthTest{
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

	negativeTest := GetPadLengthTest{
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
	positiveTests := []GetLeftRightPadLengthsTest{
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

	negativeTest := GetLeftRightPadLengthsTest{
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
	positiveTests := []GeneratePaddedStringTest{
		{
			name:          "basic generate padded string test",
			inputString:   "test",
			expected:      "#  test  #",
			char:          '#',
			numSpaces:     2,
			desiredLength: 10,
		},
	}
	for _, test := range positiveTests {
		t.Run(test.name, func(t *testing.T) {
			got, err := GeneratePaddedString(test.inputString,
				test.char, test.numSpaces, test.desiredLength)
			assertStrings(got, test.expected, t)
			assertNoError(err, t)
		})
	}
}

func TestPrintString(t *testing.T) {
	buf := &bytes.Buffer{}
	input := "this is a test string"
	err := PrintString(buf, input)
	assertNoError(err, t)
	got := buf.String()
	expected := input
	assertStrings(got, expected, t)
}

func TestGetNumCharsToPad(t *testing.T) {
	positiveTest := GetNumCharsToPadTest{
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
	negativeTest := GetNumCharsToPadTest{
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
