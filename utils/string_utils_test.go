package stringutils

import (
	"reflect"
	"testing"
)

func TestCleanInput(t *testing.T) {
	type test struct {
		name     string
		input    string
		expected []string
	}
	tests := []test{
		{
			"basic test normal spaces",
			"hello world",
			[]string{
				"hello",
				"world",
			},
		},
		{
			"test with irregular spacing",
			"    hello world     ",
			[]string{
				"hello",
				"world",
			},
		},
		{
			"test with capitalization",
			"    hElLo wOrLd     ",
			[]string{
				"hello",
				"world",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := CleanInput(test.input)
			want := test.expected
			if !reflect.DeepEqual(got, want) {
				t.Errorf("\ngot: %v\nwant: %v", got, want)
			}
		})
	}

}
