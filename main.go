package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"

	stringutils "github.com/sohWenMing/pokedex_cli/utils"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("pokedex> ")
		scanner.Scan()
		input := scanner.Text()
		command, err := getCommand(stringutils.CleanInput(input))
		if err != nil {
			continue
		}
		printCommand(command)
	}
}

func getCommand(input []string) (command string, err error) {
	if len(input) == 0 {
		return "", errors.New("string slice was empty")
	}
	return input[0], nil
}

func printCommand(command string) {
	line := fmt.Sprintf("Your command was: %s", command)
	fmt.Println(line)
}
