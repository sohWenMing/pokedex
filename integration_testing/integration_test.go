package integrationtesting

import (
	"io"
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
