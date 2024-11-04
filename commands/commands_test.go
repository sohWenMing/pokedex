package commands

import (
	"bufio"
	"bytes"
	"reflect"
	"testing"
	"time"

	cache "github.com/sohWenMing/pokedex/cache"
	testErrorHelpers "github.com/sohWenMing/pokedex/test_error_helpers"
)

var assertStrings = testErrorHelpers.AssertStrings

func TestGetCommand(t *testing.T) {

	type testStruct struct {
		name  string
		input string
		got   string
		want  string
	}

	tests := []testStruct{
		{
			"testing default",
			"should get default",
			GetCommand("should get default").name,
			defaultCommand.name,
		},
		{
			"testing help",
			"help",
			GetCommand("help").name,
			helpCommand.name,
		},
		{
			"testing exit",
			"exit",
			GetCommand("exit").name,
			exitCommand.name,
		},
		{
			"testing upper and lower",
			" eXiT      ",
			GetCommand("exit").name,
			exitCommand.name,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assertStrings(test.got, test.want, t)
		})
	}
}

func TestDefaultCallBack(t *testing.T) {
	cache := cache.NewCache(0 * time.Second)
	buf := bytes.Buffer{}
	defaultCallBack(&buf, cache)
	scanner := bufio.NewScanner(&buf)
	got := []string{}

	for scanner.Scan() {
		line := scanner.Text()
		got = append(got, line)
	}

	if !reflect.DeepEqual(got, defaultCallBackLines) {
		t.Errorf("\n got: %v\nwant: %v", got, defaultCallBackLines)
	}

}
