package main

import (
	"bufio"
	"os"

	commands "github.com/sohWenMing/pokedex/commands"
	helpers "github.com/sohWenMing/pokedex/helpers"
)

var prompt string = "pokedex"
var TrimToLower = helpers.TrimToLowerString
var runCmdCallBack = commands.ActivateCallBack

func pokeDexPrompt() {
	helpers.PrintPrompt(prompt)
}

func main() {
	helpers.PrintInitialPrompt()
	scanner := bufio.NewScanner(os.Stdin)
	var isExit bool

	pokeDexPrompt()
	//print the initial prompt

	for {
		if scanner.Scan() {
			text := TrimToLower(scanner.Text())
			isExit = runCmdCallBack(text, os.Stdout)
			if isExit {
				break
			}
			pokeDexPrompt()
		}
	}
}
