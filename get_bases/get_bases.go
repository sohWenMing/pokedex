package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	var pokemonResults AllPokemonResult
	res, err := http.Get("https://pokeapi.co/api/v2/pokemon?limit=100000")
	if err != nil {
		fmt.Println("didn't expect error, got err: ", err)
		os.Exit(1)
	}
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&pokemonResults)
	if err != nil {
		fmt.Println("didn't expect error, got err: ", err)
		os.Exit(1)
	}
	results := pokemonResults.Results
	cleanedResultsChan := make(chan CleanedPokemonResult, pokemonResults.Count)
	wg.Add(pokemonResults.Count)
	for _, result := range results {
		go func(url string) {
			defer wg.Done()
			getPokemonStats(result.URL, cleanedResultsChan)
		}(result.URL)
	}

	go func() {
		wg.Wait()
		close(cleanedResultsChan)
	}()
	min := 0
	max := 0
	totalCount := pokemonResults.Count
	countReceived := 0

	for cleanedResult := range cleanedResultsChan {
		countReceived += 1
		fmt.Println("totalCount: ", totalCount)
		fmt.Println("countReceived: ", countReceived)
		fmt.Printf("\nreceived result for %s", cleanedResult.Name)
		if cleanedResult.BaseExperience < min {
			min = cleanedResult.BaseExperience
		}
		if cleanedResult.BaseExperience > max {
			max = cleanedResult.BaseExperience
		}
	}

	fmt.Println("min: ", min)
	fmt.Println("max", max)
}

func getPokemonStats(url string, resultChan chan<- CleanedPokemonResult) {
	var pokemonResult Pokemon
	res, err := http.Get(url)
	if err != nil {
		resultChan <- CleanedPokemonResult{}
		return
	}
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&pokemonResult)
	if err != nil {
		resultChan <- CleanedPokemonResult{}
		return
	}
	cleanedResult := CleanedPokemonResult{
		pokemonResult.Name,
		pokemonResult.BaseExperience,
	}
	resultChan <- cleanedResult
	return

}
