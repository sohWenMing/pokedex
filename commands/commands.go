package commands

import (
	"fmt"
	"io"
	"log"
)

type CliCommand struct {
	Name        string
	Description string
	Callback    func(io.Writer) bool
}

var (
	mapDescriptionString  = "map: gets and displays 20 locations in the Pokemon universe"
	mapBDescriptionString = "mapb: gets and displays the previous 20 locations in the Pokemon universe"
	helpDescription       = "help: Displays a help text"
	exitDescription       = "exit: Exits the Pokedex"
	usageHeader           = "usage:"
	exitString            = "Thank you for using the Pokedex, see you next time!"
	mapString             = "This is the map function being called"
	mapBString            = "This is the mapB function being called"
)

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
	mapCommand CliCommand = CliCommand{
		"map",
		"this is to show 20 areas in the pokemon world",
		mapCallBack,
	}
	mapBCommand CliCommand = CliCommand{
		"mapb",
		"this is for showing the 20 previous areas in the pokemono world",
		mapBCallBack,
	}
)

var CommandMap = map[string]CliCommand{
	"help":  helpCommand,
	"exit":  exitCommand,
	"usage": defaultUsageCommand,
	"map":   mapCommand,
	"mapb":  mapBCommand,
}

func ActivateCallBack(text string, w io.Writer) (isExit bool) {
	command := GetCLICommand(text)
	isExit = command.Callback(w)
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
func helpCallBack(w io.Writer) (isExit bool) {
	fmt.Fprintln(w, usageHeader)
	fmt.Fprintln(w, helpDescription)
	fmt.Fprintln(w, exitDescription)
	fmt.Fprintln(w, mapDescriptionString)
	fmt.Fprintln(w, mapBDescriptionString)

	return false
}

func defaultCallBack(w io.Writer) (isExit bool) {
	fmt.Fprintln(w, helpDescription)
	fmt.Fprintln(w, exitDescription)
	return false
}

func exitCallBack(w io.Writer) (isExit bool) {
	fmt.Fprintln(w, exitString)
	return true
}

func mapCallBack(w io.Writer) (isExit bool) {
	fmt.Fprintln(w, mapString)
	return false
}

func mapBCallBack(w io.Writer) (isExit bool) {
	fmt.Fprintln(w, mapBString)
	return false
}
