package commands

import (
	"fmt"
	"log"
)

type CliCommand struct {
	Name        string
	Description string
	Callback    func() bool
}

var (
	defaultHelpString = "help: Displays a help text"
	defaultExitString = "exit: Exits the Pokedex"
	usageHelpString   = "usage: "
)

func helpCallBack() (isExit bool) {
	fmt.Println(usageHelpString)
	return false
}

func defaultCallBack() (isExit bool) {
	fmt.Println(defaultHelpString)
	fmt.Println(defaultExitString)
	return false
}

func exitCallBack() (isExit bool) {
	fmt.Println("Thank you for using the Pokedex, see you next time!")
	return true
}

var (
	helpCommand CliCommand = CliCommand{
		"help",
		"this is a command to describe help functions",
		helpCallBack,
	}
	exitCommand CliCommand = CliCommand{
		"exit",
		"this is to exit the program",
		exitCallBack,
	}
	defaultUsageCommand CliCommand = CliCommand{
		"usage",
		"this is to show the user usage suggestions if user does not key in a valid command",
		defaultCallBack,
	}
)

var CommandMap = map[string]CliCommand{
	"help":  helpCommand,
	"exit":  exitCommand,
	"usage": defaultUsageCommand,
}

func ActivateCallBack(text string) (isExit bool) {
	command := GetCLICommand(text)
	isExit = command.Callback()
	return isExit
}

func GetCLICommand(text string) (command CliCommand) {
	command, ok := CommandMap[text]
	if !ok {
		usageCommand, usageOk := CommandMap["usage"]
		if !usageOk {
			log.Fatal("program reached unrecognized command, exiting")
		}
		return usageCommand

	}
	return command
}
