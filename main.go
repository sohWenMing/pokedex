package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"

	commandcallbacks "github.com/sohWenMing/pokedex_cli/command_callbacks"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("pokedex> ")
		scanner.Scan()
		input := scanner.Text()
		err := commandcallbacks.ParseAndExecuteCommand(input)
		if err != nil {
			fmt.Println("error: ", err)
			os.Exit(1)
		}
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
