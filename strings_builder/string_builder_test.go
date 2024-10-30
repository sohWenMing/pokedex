package stringBuilder

import "testing"

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
		{
			name:          "test negative",
			inputString:   "👀",
			desiredLength: 0,
			expected:      0,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := GetPadLength(test.inputString, test.desiredLength)
			expected := test.expected
			assertValues(got, expected, t)
		})
	}
}

func TestGetLeftRightPadLengths(t *testing.T) {
	tests := []GetLeftRightPadLengthsTest{
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
		{
			name:     "testing negative",
			numToPad: -1000,
			leftPad:  0,
			rightPad: 0,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotLeft, gotRight := GetLeftRightPadLengths(test.numToPad)
			assertLeftRightPads(gotLeft, gotRight, test.leftPad, test.rightPad, t)

		})
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
