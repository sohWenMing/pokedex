package apiintegration

import (
	"fmt"
	"testing"
)

func TestGetLocation(t *testing.T) {

	config := Config{"https://pokeapi.co/api/v2/location", ""}

	res, err := GetLocation(config.Next, &config)
	if err != nil {
		t.Errorf("Got an error, %v", err)
	}
	if len(res.Results) != 20 {
		t.Errorf("\ndidn't get proper response\n%v", res)
	}
	assertStrings(res.Next, config.Next, t)
	assertStrings(res.Previous, config.Prev, t)
	fmt.Printf("\nconfig next: %s", config.Next)
	fmt.Printf("\nconfig previous: %s", config.Prev)
}

func assertStrings(got, want string, t testing.TB) {
	t.Helper()
	if got != want {
		t.Errorf("\ngot: %s\nwant: %s", got, want)
	}
}
