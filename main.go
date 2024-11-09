package main

import (
	"bufio"
	"os"
	"time"

	apiConfig "github.com/sohWenMing/pokedex/api_config"
	cache "github.com/sohWenMing/pokedex/cache"
	commands "github.com/sohWenMing/pokedex/commands"
	prompts "github.com/sohWenMing/pokedex/prompts"
)

var scanner = bufio.NewScanner(os.Stdin)
var pokeprompt = prompts.PrintPokePrompt

func main() {

	apiconfig := apiConfig.GenNewApiConfig()
	cache := cache.NewCache(10 * time.Second)

	for {
		pokeprompt(os.Stdin)
		if scanner.Scan() {
			isExit := commands.GetCommandAndFireCallBack(
				scanner.Text(),
				os.Stdout,
				cache,
				apiconfig,
			)
			if isExit {
				break
			}
		}
	}
}
