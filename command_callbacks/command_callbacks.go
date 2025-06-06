package commandcallbacks

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/sohWenMing/pokedex_cli/config"
	"github.com/sohWenMing/pokedex_cli/internal"
	structdefinitions "github.com/sohWenMing/pokedex_cli/struct_definitions"
	"github.com/sohWenMing/pokedex_cli/utils"
	stringutils "github.com/sohWenMing/pokedex_cli/utils"
)

type cliCommand struct {
	name        string
	description string
	callback    func(c *config.Config, args []string) error
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
var mapCliCommand = cliCommand{
	"map",
	"gets the next 20 location areas for the Pokemon!",
	mapCallBackfunc,
}
var mapbCliCommand = cliCommand{
	"mapb",
	"gets the previous 20 location areas for the Pokemon!",
	mapbCallBackfunc,
}
var exploreCliCommand = cliCommand{
	"explore",
	"gets list of all the pokemon that are within the are being explored",
	exploreCallbackfunc,
}
var catchCliCommand = cliCommand{
	"catch",
	"attempts to catch a pikachu - if pikachu is already caought will alert",
	catchCallbackfunc,
}
var inspectCliCommand = cliCommand{
	"catch",
	"attempts to catch a pikachu - if pikachu is already caought will alert",
	inspectCallBackfunc,
}
var pokedexCliCommand = cliCommand{
	"pokedes",
	"lists all the pokedex the user has already caught",
	pokedexCallBackfunc,
}

var callBackMap map[string]cliCommand = map[string]cliCommand{
	"exit":    exitCliCommand,
	"help":    helpCliCommand,
	"map":     mapCliCommand,
	"mapb":    mapbCliCommand,
	"explore": exploreCliCommand,
	"catch":   catchCliCommand,
	"inspect": inspectCliCommand,
	"pokedex": pokedexCliCommand,
}

func ParseAndExecuteCommand(input string, config *config.Config) error {
	commandStruct, args := ParseCommand(input)
	switch commandStruct.name {
	case "exit":
		commandStruct.callback(config, args)
		return nil
	default:
		err := commandStruct.callback(config, args)
		if err != nil {
			utilsWriteLine(config, "A problem occured when trying to get the information. Please try again")
			utilsWriteLine(config, err.Error())
		}
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

func pokedexCallBackfunc(c *config.Config, args []string) error {
	utilsWriteLine(c, fmt.Sprintf("Number of pokemon you have caught: %d", c.GetCaughtPokemon().GetNumCaught()))
	utilsWriteLine(c, "##### POKEMON #####")
	utilsWriteLine(c, c.GetCaughtPokemon().ListPokemon())
	return nil
}

func exitCallBackfunc(c *config.Config, args []string) error {
	utilsWriteLine(c, "Closing the Pokedex... Goodbye!")
	time.Sleep(1 * time.Second)
	os.Exit(0)
	return nil
}
func helpCallBackfunc(c *config.Config, args []string) error {
	utilsWriteLine(c, "Usage:")
	utilsWriteLine(c, "")
	utilsWriteLine(c, "help: Displays a help message")
	utilsWriteLine(c, "exit: Exit the Pokedex")
	utilsWriteLine(c, "map: Displays the next 20 locations in the pokemon world")
	utilsWriteLine(c, "mapb: Displays the previous 20 locations in the pokemon world")
	utilsWriteLine(c, "explore: Shows all the pokemon that are within an area")
	utilsWriteLine(c, "catch:  Attempts to catch a pokemon")
	utilsWriteLine(c, "inspect: Shows the stats of a specific pokemon")
	utilsWriteLine(c, "pokedex: Lists all the pokemon that have been caught")

	return nil
}

func mapCallBackfunc(c *config.Config, args []string) error {
	c.IncOffset()
	err := getLocationsAreas(c)
	if err != nil {
		return err
	}
	return nil
}

func mapbCallBackfunc(c *config.Config, args []string) error {
	c.DecOffSet()
	if c.GetOffSet() < 0 {
		c.ResetOffSet()
		utils.WriteLine(c.Writer, "You have reached the beginning of the map")
	}
	err := getLocationsAreas(c)
	if err != nil {
		return err
	}
	return nil
}
func exploreCallbackfunc(c *config.Config, args []string) error {
	if len(args) != 1 {
		utils.WriteLine(c.Writer, usagePrompt("explore", 1))
		return nil
	}
	urlKey := mapExploreUrl(args[0])
	err := getExploreResult(c, urlKey)
	if err != nil {
		return err
	}
	return nil
}
func catchCallbackfunc(c *config.Config, args []string) error {
	if len(args) != 1 {
		utils.WriteLine(c.Writer, usagePrompt("catch", 1))
		return nil
	}
	urlKey := mapCatchUrl(args[0])
	// first thing - work on actually calling the API
	err := callCatch(c, urlKey)
	if err != nil {
		return err
	}
	return nil
}

func inspectCallBackfunc(c *config.Config, args []string) error {
	if len(args) != 1 {
		utils.WriteLine(c.Writer, usagePrompt("inspect", 1))
		return nil
	}
	pokemonName := args[0]
	caughtPokemon, isFound := c.GetCaughtPokemon().Find(pokemonName)

	if !isFound {
		utils.WriteLine(c.Writer, fmt.Sprintf("You haven't yet caught a pokemon called %s", pokemonName))
	} else {
		utils.WriteLine(c.Writer, caughtPokemon.InspectPokemon())
	}

	return nil
}

func usagePrompt(input string, numArgs int) string {
	return fmt.Sprintf("number or arguments passed into %s must be only %d", input, numArgs)
}

func getExploreResult(c *config.Config, urlKey string) error {
	isFound, cacheEntry := getFromCache(c.GetCache(), urlKey)
	if isFound {
		utils.WriteLine(c.Writer, "Retrieving information from cache")
		pokemon := cacheEntry.WriteBufToStrings()
		writeLocsFromCachedLocations(c, pokemon)
		return nil

	} else {
		err := callExploreAPI(c, urlKey)
		if err != nil {
			return err
		}
		return nil
	}
}

func getLocationsAreas(c *config.Config) error {
	urlKey := mapLocUrl(c.GetOffSet())
	isFound, cacheEntry := getFromCache(c.GetCache(), urlKey)
	if isFound {
		utils.WriteLine(c.Writer, "Retrieving information from cache")
		locs := cacheEntry.WriteBufToStrings()
		writeLocsFromCachedLocations(c, locs)
		return nil

	} else {
		err := callLocAreas(c, urlKey)
		if err != nil {
			return err
		}
		return nil
	}
}
func callCatch(c *config.Config, url string) error {
	splitUrl := strings.Split(url, "/")
	pokemonName := splitUrl[len(splitUrl)-1]

	_, isFound := c.GetCaughtPokemon().Find(pokemonName)
	if isFound {
		utils.WriteLine(c.Writer, fmt.Sprintf("You have already caught %s!", pokemonName))
		return nil
	}

	err := callCatchAPI(c, url, pokemonName)
	if err != nil {
		return err
	}

	return nil
}

func callCatchAPI(c *config.Config, url string, pokemonName string) error {
	utils.WriteLine(c.Writer, fmt.Sprintf("Throwing a Pokeball at %s...", pokemonName))
	res, err := c.GetClient().Get(url)
	if err != nil {
		return err
	}
	if res.StatusCode == 404 {
		utils.WriteLine(c.Writer, "That pokemon doesn't exist! Try again!")
		return nil
	}
	defer res.Body.Close()
	var pokemonResult structdefinitions.Pokemon
	decoder := json.NewDecoder(res.Body)
	jsonErr := decoder.Decode(&pokemonResult)
	if jsonErr != nil {
		return err
	}
	baseExperience := pokemonResult.BaseExperience
	// stub value for now
	randResult := rand.Intn(636)
	if randResult >= baseExperience {
		utils.WriteLine(c.Writer, fmt.Sprintf("caught %s!", pokemonName))
		utils.WriteLine(c.Writer, "You may now inspect it with the inspect command")
		err = c.GetCaughtPokemon().Add(pokemonName, pokemonResult)
		if err != nil {
			c.GetCaughtPokemon().Delete(pokemonName)
		}
	} else {
		utils.WriteLine(c.Writer, fmt.Sprintf("%s escaped!", pokemonName))
	}
	return nil
}

func callLocAreas(c *config.Config, url string) error {

	res, err := c.GetClient().Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	var locAreaResult structdefinitions.LocationAreaResult
	decoder := json.NewDecoder(res.Body)
	jsonErr := decoder.Decode(&locAreaResult)
	if jsonErr != nil {
		return jsonErr
	}
	writeLocsFromLocationAreaResult(c, locAreaResult)
	err = cacheLocationAreaResults(c, url, locAreaResult)
	if err != nil {
		deleteErrorWriteFromCache(c, url)
	}
	if locAreaResult.Next == "" {
		utils.WriteLine(c.Writer, "You have reached the last page of the location areas. Will reset to start of locations on next call.")
		c.ResetOffSet()
	}
	return nil
}
func callExploreAPI(c *config.Config, url string) error {

	res, err := c.GetClient().Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	var locResult structdefinitions.LocationResult
	decoder := json.NewDecoder(res.Body)
	jsonErr := decoder.Decode(&locResult)
	if jsonErr != nil {
		return jsonErr
	}
	writePokemonFromExploreResults(c, locResult)
	err = cacheExploreResults(c, url, locResult)
	if err != nil {
		deleteErrorWriteFromCache(c, url)
	}
	return nil
}

func cacheLocationAreaResults(c *config.Config, url string, locaAreaResult structdefinitions.LocationAreaResult) error {
	names := make([]string, 0, len(locaAreaResult.Results))
	for _, result := range locaAreaResult.Results {
		name := result.Name
		names = append(names, name)
	}
	err := c.GetCache().WriteToCache(url, names)
	if err != nil {
		return err
	}
	return nil
}
func cacheExploreResults(c *config.Config, url string, locaAreaResult structdefinitions.LocationResult) error {
	names := make([]string, 0, len(locaAreaResult.PokemonEncounters))
	for _, result := range locaAreaResult.PokemonEncounters {
		name := result.Pokemon.Name
		names = append(names, name)
	}
	err := c.GetCache().WriteToCache(url, names)
	if err != nil {
		return err
	}
	return nil
}

func deleteErrorWriteFromCache(c *config.Config, key string) {
	_, ok := c.GetCache().AccessCacheMap()[key]
	if !ok {
		return
	}
	delete(c.GetCache().AccessCacheMap(), key)
}
func writePokemonFromExploreResults(c *config.Config, locationResult structdefinitions.LocationResult) {
	header := "##### Pokemon Found Start #####"
	footer := "##### Pokemon Found End   #####"
	utils.WriteLine(c.Writer, "")
	utils.WriteLine(c.Writer, header)
	utils.WriteLine(c.Writer, "")
	for _, result := range locationResult.PokemonEncounters {
		utils.WriteLine(c.Writer, result.Pokemon.Name)
	}
	utils.WriteLine(c.Writer, "")
	utils.WriteLine(c.Writer, footer)
	utils.WriteLine(c.Writer, "")
}

func writeLocsFromLocationAreaResult(c *config.Config, locationAreaResult structdefinitions.LocationAreaResult) {
	header := "##### Location Areas Start #####"
	footer := "##### Location Areas End   #####"
	utils.WriteLine(c.Writer, "")
	utils.WriteLine(c.Writer, header)
	utils.WriteLine(c.Writer, "")
	for _, locResult := range locationAreaResult.Results {
		utils.WriteLine(c.Writer, locResult.Name)
	}
	utils.WriteLine(c.Writer, "")
	utils.WriteLine(c.Writer, footer)
	utils.WriteLine(c.Writer, "")

}
func writeLocsFromCachedLocations(c *config.Config, locations []string) {
	header := "##### Location Areas Start #####"
	footer := "##### Location Areas End   #####"
	utils.WriteLine(c.Writer, "")
	utils.WriteLine(c.Writer, header)
	utils.WriteLine(c.Writer, "")
	for _, locResult := range locations {
		utils.WriteLine(c.Writer, locResult)
	}
	utils.WriteLine(c.Writer, "")
	utils.WriteLine(c.Writer, footer)
	utils.WriteLine(c.Writer, "")

}

func mapExploreUrl(param string) string {
	requesturl := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s", param)
	return requesturl
}
func mapCatchUrl(param string) string {
	requesturl := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", param)
	return requesturl
}

func mapLocUrl(offset int) string {
	requesturl := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/?offset=%d&limit=20", offset)
	return requesturl
}

func getFromCache(c *internal.Cache, key string) (isFound bool, cacheEntry internal.CacheEntry) {
	entry, ok := c.AccessCacheMap()[key]
	if !ok {
		return false, internal.CacheEntry{}
	}
	return true, entry
}

func utilsWriteLine(c *config.Config, input string) {
	utils.WriteLine(c.Writer, input)
}
