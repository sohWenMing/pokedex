package main

import (
	"bufio"
	"fmt"
	"os"

	commandcallbacks "github.com/sohWenMing/pokedex_cli/command_callbacks"
	config "github.com/sohWenMing/pokedex_cli/config"
	httputils "github.com/sohWenMing/pokedex_cli/http_utils"
)

func main() {
	config := config.InitConfig()
	config.SetClient(httputils.InitClient())

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
