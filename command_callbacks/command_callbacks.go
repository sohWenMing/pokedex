package commandcallbacks

import (
	"fmt"
	"os"
	"time"

	stringutils "github.com/sohWenMing/pokedex_cli/utils"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

var exitCliCommand = cliCommand{
	"exit",
	"user entered exit, used to exit program",
	exitCallBackfunc,
}

var helpCliCommand = cliCommand{
	"help",
	"user entered help or wrong command, show help usage",
	helpCallBackfunc,
}

var callBackMap map[string]cliCommand = map[string]cliCommand{
	"exit": exitCliCommand,
	"help": helpCliCommand,
}

func ParseAndExecuteCommand(input string) error {
	commandStruct, _ := ParseCommand(input)
	switch commandStruct.name {
	case "exit":
		commandStruct.callback()
		return nil
	default:
		commandStruct.callback()
		return nil
	}
}

func ParseCommand(input string) (commandStruct cliCommand, args []string) {
	cleanedInput := stringutils.CleanInput(input)
	if len(cleanedInput) == 0 {
		return callBackMap["help"], []string{}
	}
	commandString := cleanedInput[0]
	callBack, ok := callBackMap[commandString]
	if !ok {
		return callBackMap["help"], []string{}
	}
	return callBack, cleanedInput[1:]
}

func exitCallBackfunc() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	time.Sleep(1 * time.Second)
	os.Exit(0)
	return nil
}
func helpCallBackfunc() error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("")
	fmt.Println("help: Displays a help message")
	fmt.Println("exit: Exit the Pokedex")
	return nil
}
