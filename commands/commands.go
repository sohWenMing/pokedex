package commands

import (
	"fmt"
	"io"
	"log"

	apiconfig "github.com/sohWenMing/pokedex/api_integration"
)

type CliCommand struct {
	Name        string
	Description string
	Callback    func(w io.Writer, v *apiconfig.Config) bool
}

var (
	mapDescriptionString  = "map: gets and displays 20 locations in the Pokemon universe"
	mapBDescriptionString = "mapb: gets and displays the previous 20 locations in the Pokemon universe"
	helpDescription       = "help: Displays a help text"
	exitDescription       = "exit: Exits the Pokedex"
	usageHeader           = "usage:"
	exitString            = "Thank you for using the Pokedex, see you next time!"
)

// ########## messages to print to writer from functions ########## //
var (
	noMoreEntriesText   = "No more entries to show, going back to the beginning of entries"
	noPrevEntriesText   = "There are no previous entries to show, try using the map command"
	connectionErrorText = "There was an error connecting to the pokedex, please try again later."
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

var originalUrl = "https://pokeapi.co/api/v2/location"

var ApiConfig = apiconfig.Config{
	Next: originalUrl,
	Prev: "",
}

func ActivateCallBack(text string, w io.Writer) (isExit bool) {
	command := GetCLICommand(text)
	isExit = command.Callback(w, &ApiConfig)
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
func helpCallBack(w io.Writer, config *apiconfig.Config) (isExit bool) {
	fmt.Fprintln(w, usageHeader)
	fmt.Fprintln(w, helpDescription)
	fmt.Fprintln(w, exitDescription)
	fmt.Fprintln(w, mapDescriptionString)
	fmt.Fprintln(w, mapBDescriptionString)

	return false
}

func defaultCallBack(w io.Writer, config *apiconfig.Config) (isExit bool) {
	fmt.Fprintln(w, helpDescription)
	fmt.Fprintln(w, exitDescription)
	return false
}

func exitCallBack(w io.Writer, config *apiconfig.Config) (isExit bool) {
	fmt.Fprintln(w, exitString)
	return true
}

func mapCallBack(w io.Writer, config *apiconfig.Config) (isExit bool) {
	if config.Next == "" {
		fmt.Fprintln(w, noMoreEntriesText)
		config.Prev = ""
		config.Next = originalUrl
		return false
	}
	callApiAndSetConfig(w, config)
	return false
}

func mapBCallBack(w io.Writer, config *apiconfig.Config) (isExit bool) {
	if config.Prev == "" {
		fmt.Fprintln(w, noPrevEntriesText)
		return false
	}
	callApiAndSetConfig(w, config)

	return false
}

func callApiAndSetConfig(w io.Writer, config *apiconfig.Config) {

	response, err := apiconfig.GetLocation(config.Next, config)
	if err != nil {
		fmt.Fprintln(w, connectionErrorText)
		return
	}
	for i, result := range response.Results {
		fmt.Fprintf(w, "Result %d: %v\n", i+1, result.Name)
	}
}
