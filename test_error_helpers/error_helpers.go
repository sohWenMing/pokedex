package testErrorHelpers

import "testing"

func AssertStrings(got, want string, t testing.TB) {
	t.Helper()
	if got != want {
		t.Errorf("\ngot: %s\nwant: %s", got, want)
	}
}

func AssertVals(got, want int, t testing.TB) {
	t.Helper()
	if got != want {
		t.Errorf("\ngot: %d\nwant: %d", got, want)
	}
}
