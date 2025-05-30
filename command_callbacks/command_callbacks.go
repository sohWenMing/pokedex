package commandcallbacks

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/sohWenMing/pokedex_cli/config"
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

var callBackMap map[string]cliCommand = map[string]cliCommand{
	"exit": exitCliCommand,
	"help": helpCliCommand,
	"map":  mapCliCommand,
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
	requesturl := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/?offset=%d&limit=20", c.GetOffSet())
	header := "##### Location Areas Start #####"
	footer := "##### Location Areas End   #####"

	utils.WriteLine(c.Writer, header)
	res, err := c.GetClient().Get(requesturl)
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
	for _, loc_area := range locAreaResult.Results {
		utils.WriteLine(c.Writer, loc_area.Name)
	}
	if locAreaResult.Next == "" {
		utils.WriteLine(c.Writer, "You have reached the last page of the location areas. Will reset to start of locations on next call.")
		c.ResetOffSet()
	}
	utils.WriteLine(c.Writer, footer)

	return nil
}
