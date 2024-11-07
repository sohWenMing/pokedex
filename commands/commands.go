package commands

import (
	"fmt"
	"io"
	"log"
	"strings"

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

var (
	helpCommand    = command{"help", helpCallBack}
	exitCommand    = command{"exit", exitCallBack}
	defaultCommand = command{"default", defaultCallBack}
	mapCommand     = command{"map", mapCallBack}
	mapBCommand    = command{"mapb", mapBCallBack}
)

var commandMap = map[string]command{
	"help":    helpCommand,
	"exit":    exitCommand,
	"default": defaultCommand,
	"map":     mapCommand,
	"mapb":    mapBCommand,
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
