package commands

import (
	"fmt"
	"io"
	"strings"
	"time"

	apiConfig "github.com/sohWenMing/pokedex/api_config"
	cache "github.com/sohWenMing/pokedex/cache"
	prompts "github.com/sohWenMing/pokedex/prompts"
)

type command struct {
	name     string
	Callback func(io.Writer, *cache.Cache, *apiConfig.ApiConfig) (isExit bool)
}

var defaultCallBackLines = []string{
	"usage:",
	"help: shows list of commands",
	"map: gets next 20 locations",
	"mapb: gets previous 20 locations",
	"exit: exit",
}

var helpCallBackLines = []string{
	"here are some ways you can use the pokedex:",
	"",
	"help:",
	"will show you this helpful help view",
	"",
	"map:",
	"gets the next available locations.",
	"If no more next locations are available, will set configuration so that the when map is next called, will get the first 20 locations",
	"",
	"mapb:",
	"gets the previous available locations. If no more previous locations are available, will alert you.",
	"",
	"exit:",
	"exits the pokedex",
}

func defaultCallBack(w io.Writer, c *cache.Cache, a *apiConfig.ApiConfig) (isExit bool) {
	printLines(w, defaultCallBackLines)
	return false
}

func exitCallBack(w io.Writer, c *cache.Cache, a *apiConfig.ApiConfig) (isExit bool) {
	prompts.PrintExitPrompt(w)
	return true
}

func helpCallBack(w io.Writer, c *cache.Cache, a *apiConfig.ApiConfig) (isExit bool) {
	printLines(w, helpCallBackLines)
	return false
}

func mapCallBack(w io.Writer, c *cache.Cache, a *apiConfig.ApiConfig) (isExit bool) {
	fmt.Printf("getting information...")
	time.Sleep(2 * time.Second)
	return false
}

var (
	helpCommand    = command{"help", helpCallBack}
	exitCommand    = command{"exit", exitCallBack}
	defaultCommand = command{"default", defaultCallBack}
	mapCommand     = command{"map", mapCallBack}
)

var commandMap = map[string]command{
	"help":    helpCommand,
	"exit":    exitCommand,
	"default": defaultCommand,
	"map":     mapCommand,
}

func GetCommand(input string) command {
	formattedInput := strings.ToLower(strings.TrimSpace(input))

	_, ok := commandMap[formattedInput]
	if !ok {
		return commandMap["default"]
	}
	return commandMap[formattedInput]
}

func printLines(w io.Writer, strings []string) {
	for _, line := range strings {
		fmt.Fprintln(w, line)
	}
}
