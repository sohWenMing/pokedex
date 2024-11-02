package stringbuilder

import "fmt"

func buildPrompt(promptText, promptChar string) string {
	returnedString := fmt.Sprintf("%s%s ", promptText, promptChar)
	return returnedString
}

var PokePrompt = buildPrompt("pokedex", ">")
