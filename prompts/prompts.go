package prompts

import (
	"fmt"
	"io"

	colors_package "github.com/sohWenMing/pokedex/color_package"
	stringbuilder "github.com/sohWenMing/pokedex/string_builder"
)

var writeBlue = colors_package.WriteBlue
var writeRed = colors_package.WriteRed
var pokePrompt = stringbuilder.PokePrompt
var exitPrompt = fmt.Sprintln("Thank you for using the pokedex. See you soon!")

func PrintPokePrompt(w io.Writer) {
	writeBlue(w, pokePrompt)
}

func PrintExitPrompt(w io.Writer) {
	writeRed(w, exitPrompt)
}
