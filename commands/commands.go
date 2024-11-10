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
	Callback func(w io.Writer, c *cache.Cache, a *apiConfig.ApiConfig, explore_location string) (isExit bool)
}

const (
	exploreUsageString string = `usage - "explore <enter location here>"`
	exploreErrorString string = "location entered could not be found. Please try again."
	exploringString    string = "exploring for pokemon ... "
	numPokemonsString  string = "number of pokemons in explored location: <num>"
)

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

func defaultCallBack(w io.Writer, c *cache.Cache, a *apiConfig.ApiConfig, explore_location string) (isExit bool) {
	printLines(w, defaultCallBackLines)
	return false
}

func exitCallBack(w io.Writer, c *cache.Cache, a *apiConfig.ApiConfig, explore_location string) (isExit bool) {
	prompts.PrintExitPrompt(w)
	return true
}

func helpCallBack(w io.Writer, c *cache.Cache, a *apiConfig.ApiConfig, explore_location string) (isExit bool) {
	printLines(w, helpCallBackLines)
	return false
}

func mapCallBack(w io.Writer, c *cache.Cache, a *apiConfig.ApiConfig, explore_location string) (isExit bool) {
	fmt.Fprintf(w, "getting information...\n")

	// check the information in the cache
	urlToCall := a.GetNext()
	mapOutputVals(urlToCall, true, w, c, a)

	return false
}

func mapBCallBack(w io.Writer, c *cache.Cache, a *apiConfig.ApiConfig, explore_location string) (isExit bool) {
	fmt.Fprintf(w, "getting information...\n")

	urlToCall := a.GetPrev()
	mapOutputVals(urlToCall, false, w, c, a)
	return false
}

func exploreCallBack(w io.Writer, c *cache.Cache, a *apiConfig.ApiConfig, explore_location string) (isExit bool) {

	// if explore was entered, but no location was entered, then show usage message
	if explore_location == "" {
		fmt.Fprint(w, exploreUsageString)
		return false
	}
	// show that the program is searching for pokemon
	fmt.Fprintln(w, exploringString)

	// get the pokemons by calling the PokeAPI, checking for returned errors
	pokemons, err := apiConfig.CallExploreUrl(explore_location)

	// if error is returned, show error message
	if err != nil {
		fmt.Fprintln(w, exploreErrorString)
		return false
	}
	// first print the number of pokemon found

	numPokemons := len(pokemons)
	numPokemonsString := helpers.StringReplace(numPokemonsString, "<num>", fmt.Sprintf("%d", numPokemons))
	fmt.Fprintln(w, numPokemonsString)
	for _, pokemon := range pokemons {

		fmt.Fprintln(w, pokemon.Name)
	}
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

	command, location := GetCommand(inputString)
	isExit = command.Callback(w, c, a, location)
	return isExit
}

func GetCommand(input string) (command command, location string) {
	formattedInput := helpers.ToLowerAndTrim(input)
	commandString, location := getCommandStringAndLocation(formattedInput)
	_, ok := commandMap[commandString]
	if !ok {
		return commandMap["default"], ""
	}
	return commandMap[commandString], location
}

func printLines(w io.Writer, strings []string) {
	for _, line := range strings {
		fmt.Fprintln(w, line)
	}
}

func mapOutputVals(urlToCall string, isNext bool, w io.Writer, c *cache.Cache, a *apiConfig.ApiConfig) {

	//if the values can be found in the cache, the write to the writer and exit early
	next, prev, firstCacheCallValues, firstCacheCallErr := c.GetFromCache(urlToCall)
	if firstCacheCallErr == nil {
		for _, value := range firstCacheCallValues {
			fmt.Fprintln(w, value.Name)
		}
		a.SetConfig(next, prev)
		return
	}

	next, prev, values, callErr := a.CallMapUrl(isNext)
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

func getCommandStringAndLocation(formattedInput string) (commandString, location string) {
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

	locationLoweredAndTrimmed := helpers.ToLowerAndTrim(splitStrings[1])
	location = helpers.ReplaceSpaceWithChar(locationLoweredAndTrimmed, "-")
	return commandString, location
}
