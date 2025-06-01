package integrationtesting

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"
	"time"

	commandcallbacks "github.com/sohWenMing/pokedex_cli/command_callbacks"
	"github.com/sohWenMing/pokedex_cli/config"
	httputils "github.com/sohWenMing/pokedex_cli/http_utils"
)

func TestResetOffset(t *testing.T) {
	config, err := config.InitConfig(io.Discard, 1*time.Second, 1*time.Second)
	if err != nil {
		t.Errorf("\ndidn't expect error, got %v", err)
	}
	config.SetClient(httputils.InitClient())

	for config.GetOffSet() <= 10000 {
		config.IncOffset()
	}
	err = commandcallbacks.ParseAndExecuteCommand("map", &config)
	if err != nil {
		t.Errorf("\ndidn't expect error, got %v", err)
	}
	got := config.GetOffSet()
	want := 0
	if got != want {
		t.Errorf("\ngot: %d\nwant: %d", got, want)
	}
}
func TestExplore(t *testing.T) {
	areas := []string{
		// "canalave-city-area",
		// "eterna-city-area",
		// "pastoria-city-area",
		// "sunyshore-city-area",
		// "sinnoh-pokemon-league-area",
		// "oreburgh-mine-1f",
		// "oreburgh-mine-b1f",
		// "valley-windworks-area",
		// "eterna-forest-area",
		// "fuego-ironworks-area",
		// "mt-coronet-1f-route-207",
		// "mt-coronet-2f",
		// "mt-coronet-3f",
		// "mt-coronet-exterior-snowfall",
		// "mt-coronet-exterior-blizzard",
		// "mt-coronet-4f",
		// "mt-coronet-4f-small-room",
		// "mt-coronet-5f",
		// "mt-coronet-6f",
		"mt-coronet-1f-from-exterior",
	}
	buf := &bytes.Buffer{}
	expectedPokemon := []string{
		"clefairy",
		"golbat",
		"machoke",
		"graveler",
		"nosepass",
		"medicham",
		"lunatone",
		"solrock",
		"chingling",
		"bronzong",
	}

	config, err := config.InitConfig(buf, 1*time.Second, 1*time.Second)
	if err != nil {
		t.Errorf("\ndidn't expect error, got %v", err)
	}
	config.SetClient(httputils.InitClient())

	for _, area := range areas {
		command := fmt.Sprintf("explore %s", area)

		err = commandcallbacks.ParseAndExecuteCommand(command, &config)
		if err != nil {
			t.Errorf("\ndidn't expect error, got %v", err)
		}
	}
	for _, pokemon := range expectedPokemon {
		if !strings.Contains(buf.String(), pokemon) {
			t.Errorf("\n%s was expected but not found", pokemon)
		}
	}
}
