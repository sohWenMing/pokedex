package stringbuilder

import (
	"testing"

	testErrorHelpers "github.com/sohWenMing/pokedex/test_error_helpers"
)

var assertStrings = testErrorHelpers.AssertStrings

func TestBuildPrompt(t *testing.T) {
	type testStruct struct {
		name       string
		promptText string
		promptChar string
		want       string
	}

	tests := []testStruct{
		{
			"basic test",
			"test",
			">",
			"test> ",
		},
		{
			"basic test",
			"test",
			"💩",
			"test💩 ",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := buildPrompt(test.promptText, test.promptChar)
			want := test.want
			if got != want {
				assertStrings(got, want, t)
			}
		})
	}
}
