package helpers

import (
	"testing"

	test_error_helpers "github.com/sohWenMing/pokedex/test_error_helpers"
)

func TestToLowerAndTrim(t *testing.T) {
	got := ToLowerAndTrim(" hEre   ")
	want := "here"
	test_error_helpers.AssertStrings(got, want, t)
}

func TestReplaceSpaceWithChar(t *testing.T) {
	got := ReplaceSpaceWithChar("this is a test", "-")
	want := "this-is-a-test"
	test_error_helpers.AssertStrings(got, want, t)
}
