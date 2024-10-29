package main

import (
	"bufio"
	"fmt"
	"os"

	helpers "github.com/sohWenMing/pokedex/helpers"
	// commands "github.com/sohWenMing/pokedex/commands"
)

var prompt string = "pokedex"
var TrimToLower = helpers.TrimToLowerString

func pokeDexPrompt() {
	helpers.PrintPrompt(prompt)
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	pokeDexPrompt()
	//print the initial prompt

	for {
		if scanner.Scan() {
			text := TrimToLower(scanner.Text())
			fmt.Printf("%q", text)
			pokeDexPrompt()
		}
	}

}
