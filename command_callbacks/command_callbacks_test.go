package commandcallbacks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	structdefinitions "github.com/sohWenMing/pokedex_cli/struct_definitions"
)

func TestNullNext(t *testing.T) {
	requesturl := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/?offset=10000&limit=20")
	res, err := http.Get(requesturl)
	if err != nil {
		t.Errorf("didnt expect error: got %v", err)
	}
	defer res.Body.Close()
	var locAreaResult structdefinitions.LocationAreaResult
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&locAreaResult)
	if err != nil {
		t.Errorf("didnt expect error: got %v", err)
	}
	fmt.Printf("\nvalue of next: %v", locAreaResult.Next)
	fmt.Printf("\nIs next blank: %v", locAreaResult.Next == "")
}

func TestNonExistingPokemon(t *testing.T) {
	pokemon := "thishouldfail"
	url := mapCatchUrl(pokemon)
	res, err := http.Get(url)
	if err == nil {
		fmt.Print(err)
	}
	fmt.Println("status code: ", res.StatusCode)
}
