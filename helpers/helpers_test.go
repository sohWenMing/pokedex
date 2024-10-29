package helpers

import (
	"testing"
)

type stringTest struct {
	name     string
	input    string
	expected string
}

var expectedText = "test text"

func TestTrimToLower(t *testing.T) {
	tests := []stringTest{
		{
			name:     "testing spacing",
			input:    "            test text            ",
			expected: expectedText,
		},
		{
			name:     "testing to lower",
			input:    "    tEsT TEXT",
			expected: expectedText,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := TrimToLowerString(test.input)
			want := test.expected
			assertCorrectValue(got, want, t)
		})
	}
}

func assertCorrectValue(got, want string, t testing.TB) {
	t.Helper()
	if got != want {
		t.Errorf("\ngot: %s\nwant: %s\n", got, want)
	}
}
