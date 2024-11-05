package commands

import (
	"bufio"
	"bytes"
	"reflect"
	"testing"
	"time"

	apiCfg "github.com/sohWenMing/pokedex/api_config"
	cache "github.com/sohWenMing/pokedex/cache"
	colors_package "github.com/sohWenMing/pokedex/color_package"
	prompts "github.com/sohWenMing/pokedex/prompts"
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

func TestDefaultAndExitCallBack(t *testing.T) {
	cache := cache.NewCache(0 * time.Second)
	buf := bytes.Buffer{}
	apiconfig := apiCfg.GenNewApiConfig()

	//testing the printout from the default callback
	defaultCallBack(&buf, cache, apiconfig)
	scanner := bufio.NewScanner(&buf)
	got := []string{}

	for scanner.Scan() {
		line := scanner.Text()
		got = append(got, line)
	}

	if !reflect.DeepEqual(got, defaultCallBackLines) {
		t.Errorf("\n got: %v\nwant: %v", got, defaultCallBackLines)
	}

	//testing the printout from the exit callBack
	exitBuf := bytes.Buffer{}
	wantBuf := bytes.Buffer{}
	exitCallBack(&exitBuf, cache, apiconfig)

	gotExitPrompt := exitBuf.String()

	colors_package.WriteRed(&wantBuf, prompts.GetExitPrompt())
	wantExitPrompt := wantBuf.String()
	assertStrings(gotExitPrompt, wantExitPrompt, t)

}

func TestMapCallBack(t *testing.T) {
	cache := cache.NewCache(1 * time.Second)
	apiconfig := apiCfg.GenNewApiConfig()

	urlToCall := apiconfig.GetNext()

	//first assert that there are no values tied to the URL
	cacheCallValues, cacheCallErr := cache.GetFromCache(urlToCall)
	testErrorHelpers.AssertError(cacheCallErr, t)
	testErrorHelpers.AssertReflectDeepEqual(cacheCallValues, []apiCfg.MapValue{}, t)

	buf := bytes.Buffer{}
	mapCallBack(&buf, cache, apiconfig)
	//function should call the API, and write returned values to cache

	scanner := bufio.NewScanner(&buf)
	stringsFromCache := []string{}
	for scanner.Scan() {
		text := scanner.Text()
		stringsFromCache = append(stringsFromCache, text)
	}
	testErrorHelpers.AssertVals(len(stringsFromCache), 20, t)

}
