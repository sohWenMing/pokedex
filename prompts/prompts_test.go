package prompts

import (
	"bytes"
	"testing"

	colors_package "github.com/sohWenMing/pokedex/color_package"
	test_error_helpers "github.com/sohWenMing/pokedex/test_error_helpers"
)

var assertStrings = test_error_helpers.AssertStrings

func TestPokePrompt(t *testing.T) {
	buf := bytes.Buffer{}
	PrintPokePrompt(&buf)
	got := buf.String()
	want := colors_package.ColorBlue.Sprintf(pokePrompt)
	assertStrings(got, want, t)
}

func TestExitPrompt(t *testing.T) {
	buf := bytes.Buffer{}
	PrintExitPrompt(&buf)
	got := buf.String()
	want := colors_package.ColorRed.Sprintf(exitPrompt)
	assertStrings(got, want, t)
}
