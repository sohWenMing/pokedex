package main

import (
	"bufio"
	"os"

	commands "github.com/sohWenMing/pokedex/commands"
	prompts "github.com/sohWenMing/pokedex/prompts"
)

var scanner = bufio.NewScanner(os.Stdin)
var getCommand = commands.GetCommand
var pokeprompt = prompts.PrintPokePrompt

func main() {
	for {
		pokeprompt(os.Stdin)
		if scanner.Scan() {
			isExit := getCommand(scanner.Text()).Callback(os.Stdout)
			if isExit {
				break
			}
		}

	}
}
