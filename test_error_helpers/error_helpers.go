package testErrorHelpers

import (
	"reflect"
	"testing"
)

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

func AssertReflectDeepEqual(got, want any, t testing.TB) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("\ngot: %v\nwant: %v", got, want)
	}
}

func AssertNoError(err error, t testing.TB) {
	t.Helper()
	if err != nil {
		t.Errorf("\ndidn't expect error, got %s", err.Error())
	}
}

func AssertError(err error, t testing.TB) {
	t.Helper()
	if err == nil {
		t.Errorf("\nexpected error, didn't get one")
	}
}
