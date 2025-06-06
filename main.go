package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	commandcallbacks "github.com/sohWenMing/pokedex_cli/command_callbacks"
	config "github.com/sohWenMing/pokedex_cli/config"
	httputils "github.com/sohWenMing/pokedex_cli/http_utils"
)

func main() {
	config, err := config.InitConfig(os.Stdout, 1*time.Minute, 5*time.Minute)
	if err != nil {
		fmt.Println("an error occured: exiting program")
		os.Exit(1)
	}
	config.SetClient(httputils.InitClient())
	fmt.Println("Welcome to the Pokedex! Enter help to get instructions.")

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("pokedex> ")
		scanner.Scan()
		input := scanner.Text()
		err := commandcallbacks.ParseAndExecuteCommand(input, &config)
		if err != nil {
			fmt.Println("error: ", err)
			os.Exit(1)
		}
	}
}
