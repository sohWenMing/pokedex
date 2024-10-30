package helpers

import (
	"fmt"
	"log"
	"os"
	"strings"

	stringBuilder "github.com/sohWenMing/pokedex/strings_builder"
)

func PrintInitialPrompt() {
	for i := 0; i < 5; i++ {
		fmt.Print("\n")
	}
	paddedString, err := stringBuilder.GeneratePaddedString(
		"", '#', 0, 50,
	)
	if err != nil {
		log.Fatal("an error occured when generating paddded string", err)
	}
	printErr := stringBuilder.PrintString(os.Stdout, paddedString)
	if printErr != nil {
		log.Fatal("an error occured during printing of initial prompt", printErr)
	}
	fmt.Print("\n")

}

func PrintPrompt(input string) {
	fmt.Printf("%s> ", input)
}

func TrimToLowerString(input string) string {
	input = strings.TrimSpace(input)
	input = strings.ToLower(input)
	return input
}
