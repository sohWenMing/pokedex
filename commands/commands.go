package commands

import (
	"fmt"
	"io"
	"strings"

	prompts "github.com/sohWenMing/pokedex/prompts"
)

type command struct {
	name     string
	printOut string
	Callback func(io.Writer) (isExit bool)
}

var defaultCallBackLines = []string{
	"usage:",
	"help: shows list of commands",
	"map: gets next 20 locations",
	"mapb: gets previous 20 locations",
	"exit: exit",
}

func defaultCallBack(w io.Writer) (isExit bool) {
	for _, line := range defaultCallBackLines {
		fmt.Fprintln(w, line)
	}
	return false
}

func exitCallBack(w io.Writer) (isExit bool) {
	prompts.PrintExitPrompt(w)
	return true
}

func helpCallBack(w io.Writer) (isExit bool) {
	return false
}

var (
	helpCommand    = command{"help", "help printout", helpCallBack}
	exitCommand    = command{"exit", "exit printout", exitCallBack}
	defaultCommand = command{"default", "default printout", defaultCallBack}
)

var commandMap = map[string]command{
	"help":    helpCommand,
	"exit":    exitCommand,
	"default": defaultCommand,
}

func GetCommand(input string) command {
	formattedInput := strings.ToLower(strings.TrimSpace(input))

	_, ok := commandMap[formattedInput]
	if !ok {
		return commandMap["default"]
	}
	return commandMap[formattedInput]
}
