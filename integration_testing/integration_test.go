package integrationtesting

import (
	"io"
	"testing"

	commandcallbacks "github.com/sohWenMing/pokedex_cli/command_callbacks"
	"github.com/sohWenMing/pokedex_cli/config"
	httputils "github.com/sohWenMing/pokedex_cli/http_utils"
)

func TestResetOffset(t *testing.T) {
	config := config.InitConfig(io.Discard)
	config.SetClient(httputils.InitClient())

	for config.GetOffSet() <= 10000 {
		config.IncOffset()
	}

	err := commandcallbacks.ParseAndExecuteCommand("map", &config)
	if err != nil {
		t.Errorf("\ndidn't expect error, got %v", err)
	}
	got := config.GetOffSet()
	want := 0
	if got != want {
		t.Errorf("\ngot: %d\nwant: %d", got, want)
	}
}
