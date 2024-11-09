package commands

import (
	"fmt"
	"io"
	"log"
	"strings"

	apiConfig "github.com/sohWenMing/pokedex/api_config"
	cache "github.com/sohWenMing/pokedex/cache"
	helpers "github.com/sohWenMing/pokedex/helpers"
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
	fmt.Fprintf(w, "getting information...\n")

	// check the information in the cache
	urlToCall := a.GetNext()
	outputValues(urlToCall, true, w, c, a)

	return false
}

func mapBCallBack(w io.Writer, c *cache.Cache, a *apiConfig.ApiConfig) (isExit bool) {
	fmt.Fprintf(w, "getting information...\n")

	urlToCall := a.GetPrev()
	outputValues(urlToCall, false, w, c, a)
	return false
}

func exploreCallBack(w io.Writer, c *cache.Cache, a *apiConfig.ApiConfig) (isExit bool) {
	fmt.Println("explore callback was called")
	return false
}

var (
	helpCommand    = command{"help", helpCallBack}
	exitCommand    = command{"exit", exitCallBack}
	defaultCommand = command{"default", defaultCallBack}
	mapCommand     = command{"map", mapCallBack}
	mapBCommand    = command{"mapb", mapBCallBack}
	exploreCommand = command{"explore", exploreCallBack}
)

var commandMap = map[string]command{
	"help":    helpCommand,
	"exit":    exitCommand,
	"default": defaultCommand,
	"map":     mapCommand,
	"mapb":    mapBCommand,
	"explore": exploreCommand,
}

func GetCommandAndFireCallBack(
	inputString string,
	w io.Writer,
	c *cache.Cache,
	a *apiConfig.ApiConfig) (isExit bool) {

	command := GetCommand(inputString)
	isExit = command.Callback(w, c, a)
	return isExit
}

func GetCommand(input string) command {
	formattedInput := helpers.ToLowerAndTrim(input)
	commandString, _ := getCommandString(formattedInput)
	_, ok := commandMap[commandString]
	if !ok {
		return commandMap["default"]
	}
	return commandMap[commandString]
}

func printLines(w io.Writer, strings []string) {
	for _, line := range strings {
		fmt.Fprintln(w, line)
	}
}

func outputValues(urlToCall string, isNext bool, w io.Writer, c *cache.Cache, a *apiConfig.ApiConfig) {

	//if the values can be found in the cache, the write to the writer and exit early
	next, prev, firstCacheCallValues, firstCacheCallErr := c.GetFromCache(urlToCall)
	if firstCacheCallErr == nil {
		for _, value := range firstCacheCallValues {
			fmt.Fprintln(w, value.Name)
		}
		a.SetConfig(next, prev)
		return
	}

	next, prev, values, callErr := a.CallUrl(isNext)
	// if error occured during the API call dur to connection, print the error and return
	if callErr != nil {
		fmt.Fprintln(w, callErr.Error())
		return
	}
	a.SetConfig(next, prev)

	c.WriteToCache(urlToCall, next, prev, values)

	_, _, secondCacheCallValues, secondCacheCallErr := c.GetFromCache(urlToCall)

	// at this point, if there are still values in the cache, program should crash and investigation should happen
	if secondCacheCallErr != nil {
		log.Fatal("problem with caching mechanism")
	}

	for _, value := range secondCacheCallValues {
		fmt.Fprintln(w, value.Name)
	}
}

func checkIsExploreCommand(input string) (isExplore bool, isHasLocation bool) {
	if input == "explore" {
		isExplore = true
		isHasLocation = false
		return
	}
	if strings.HasPrefix(input, "explore ") {
		isExplore = true
		isHasLocation = true
		return
	}
	isExplore = false
	isHasLocation = false
	return
}

func getCommandString(formattedInput string) (commandString, location string) {
	isExplore, isHasLocation := checkIsExploreCommand(formattedInput)
	if !isExplore {
		commandString = formattedInput
		location = ""
		return
	}
	if isExplore && !isHasLocation {
		commandString = "explore"
		location = ""
		return
	}

	splitStrings := strings.SplitAfterN(formattedInput, " ", 2)

	commandString = helpers.ToLowerAndTrim(splitStrings[0])

	location = helpers.ToLowerAndTrim(splitStrings[1])
	return
}
