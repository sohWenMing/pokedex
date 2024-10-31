package helpers

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	stringBuilder "github.com/sohWenMing/pokedex/strings_builder"
)

func PrintInitialPrompt() {
	for i := 0; i < 2; i++ {
		fmt.Print("\n")
	}
	//print two new lines
	promptStrings := initPromptStrings()
	printFromStringSlice(promptStrings)
	for i := 0; i < 2; i++ {
		fmt.Print("\n")
	}
}

func PrintPrompt(input string) {
	// fmt.Printf("%s> ", input)
	color.New(color.FgBlue).Add(color.Underline).Fprintf(os.Stdout, "%s>", input)
	fmt.Print(" ")
}

func TrimToLowerString(input string) string {
	input = strings.TrimSpace(input)
	input = strings.ToLower(input)
	return input
}

func initPromptStrings() []string {

	type paddedStringStruct struct {
		text     string
		numSpace int
		padChar  rune
	}

	paddedStringStructs := []paddedStringStruct{
		{
			text:     "",
			numSpace: 0,
			padChar:  '#',
		},
		{
			text:     "",
			numSpace: 0,
			padChar:  '*',
		}, {
			text:     "Welcome",
			numSpace: 4,
			padChar:  ' ',
		}, {
			text:     "To",
			numSpace: 4,
			padChar:  ' ',
		}, {
			text:     "Nindgabeet's AMAAAZING",
			numSpace: 4,
			padChar:  ' ',
		}, {
			text:     "POKEDEX!!!",
			numSpace: 4,
			padChar:  ' ',
		}, {
			text:     "",
			numSpace: 0,
			padChar:  '*',
		},
		{
			text:     "",
			numSpace: 0,
			padChar:  '#',
		},
	}

	returnedStrings := []string{}

	for _, padStruct := range paddedStringStructs {
		text, numSpace, padChar := padStruct.text, padStruct.numSpace, padStruct.padChar
		blue := color.New(color.FgBlue).Add(color.Bold).SprintFunc()
		coloredText := fmt.Sprint(blue(text))

		paddedString, err := stringBuilder.GeneratePaddedString(
			coloredText, padChar, numSpace, 50,
		)
		if err != nil {
			log.Fatal("an error occured when generating padded string", err)
		}
		returnedStrings = append(returnedStrings, fmt.Sprintf("%s\n", paddedString))
	}

	return returnedStrings
}

func printFromStringSlice(inputSlice []string) {
	for _, string := range inputSlice {
		printErr := stringBuilder.PrintString(os.Stdout, string)
		if printErr != nil {
			log.Fatal("an error occured when printing strings from slices", printErr)
		}
		time.Sleep(time.Millisecond * 100)
	}
}
