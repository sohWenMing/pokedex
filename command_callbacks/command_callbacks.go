package commandcallbacks

import (
	"encoding/json"
	"fmt"
	"os"
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
	callback    func(*config.Config) error
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

var callBackMap map[string]cliCommand = map[string]cliCommand{
	"exit": exitCliCommand,
	"help": helpCliCommand,
	"map":  mapCliCommand,
	"mapb": mapbCliCommand,
}

func ParseAndExecuteCommand(input string, config *config.Config) error {
	commandStruct, _ := ParseCommand(input)
	switch commandStruct.name {
	case "exit":
		commandStruct.callback(config)
		return nil
	case "map":
		err := commandStruct.callback(config)
		if err != nil {
			fmt.Println("A problem occured when trying to get the information. Please try again")
			fmt.Println(err)
		}
		return nil

	default:
		commandStruct.callback(config)
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

func exitCallBackfunc(*config.Config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	time.Sleep(1 * time.Second)
	os.Exit(0)
	return nil
}
func helpCallBackfunc(*config.Config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("")
	fmt.Println("help: Displays a help message")
	fmt.Println("exit: Exit the Pokedex")
	return nil
}
func mapCallBackfunc(c *config.Config) error {
	c.IncOffset()
	err := getLocationsAreas(c)
	if err != nil {
		return err
	}
	return nil
}

func mapbCallBackfunc(c *config.Config) error {
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

func deleteErrorWriteFromCache(c *config.Config, key string) {
	_, ok := c.GetCache().AccessCacheMap()[key]
	if !ok {
		return
	}
	delete(c.GetCache().AccessCacheMap(), key)
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
