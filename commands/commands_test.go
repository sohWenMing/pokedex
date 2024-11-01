package commands

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"reflect"
	"strings"
	"testing"

	apiIntegration "github.com/sohWenMing/pokedex/api_integration"
	helpers "github.com/sohWenMing/pokedex/helpers"
)

type testStruct struct {
	name                   string
	input                  string
	expectedStringsInprint []string
	expected               CliCommand
}

var testApiConfig = apiIntegration.Config{
	Next: originalUrl,
	Prev: "",
}

var buf = &bytes.Buffer{}

var replaceNewLines = helpers.ReplaceNewLines
var trimString = helpers.TrimString

var (
	offsetTwentyUrl = "https://pokeapi.co/api/v2/location?offset=20&limit=20"
	offsetFortyUrl  = "https://pokeapi.co/api/v2/location?offset=40&limit=20"
	offsetZeroUrl   = "https://pokeapi.co/api/v2/location?offset=0&limit=20"
)

func TestGetCLICommand(t *testing.T) {

	tests := []testStruct{
		{
			name:  "testing help command",
			input: "help",

			expected:               helpCommand,
			expectedStringsInprint: []string{usageHeader, helpDescription, exitDescription, mapDescriptionString, mapBDescriptionString},
		},
		{
			name:                   "testing no command",
			input:                  "erroneous text",
			expected:               defaultUsageCommand,
			expectedStringsInprint: []string{helpDescription, exitDescription},
		},
		{
			name:                   "testing exit command",
			input:                  "exit",
			expected:               exitCommand,
			expectedStringsInprint: []string{exitString},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := GetCLICommand(test.input)
			assertValues(got, test.expected, t)
			got.Callback(buf, &testApiConfig)
			assertStringContains(
				buf.String(), generateExpectedPrint(test.expectedStringsInprint), t)
		})
	}
	resetBuffer()
}

func TestMapCallBack(t *testing.T) {
	t.Run("testing that api config is being set correctly by mapCallBack function", func(t *testing.T) {

		// reset the test api config
		resetTestApiConfig()
		assertStrings(testApiConfig.Next, originalUrl, t)
		assertStrings(testApiConfig.Prev, "", t)

		//make the first call to the API, config should set a next
		mapCallBack(buf, &testApiConfig)
		assertStrings(testApiConfig.Next, offsetTwentyUrl, t)
		assertStrings(testApiConfig.Prev, "", t)

		//make the second call, config should update next and prev
		mapCallBack(buf, &testApiConfig)
		assertStrings(testApiConfig.Next, offsetFortyUrl, t)
		assertStrings(testApiConfig.Prev, offsetZeroUrl, t)

		resetBuffer()
		//clear all the writes to the test buffer made by the prev api calls
		testApiConfig.Next = ""
		mapCallBack(buf, &testApiConfig)
		assertStrings(trimString(replaceNewLines(buf.String())), noMoreEntriesText, t)
		resetBuffer()
	})

	t.Run("testing that I am getting the correct number of responses", func(t *testing.T) {
		resetTestApiConfig()
		mapCallBack(buf, &testApiConfig)
		var lines []string
		scanner := bufio.NewScanner(buf)
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			t.Errorf("error occured during scanning: %s", err.Error())
		}
		if len(lines) != 20 {
			t.Errorf("expected 20 lines, got %d lines", len(lines))
		}

	})

}

func TestMapBCallBack(t *testing.T) {
	t.Run("testing that api config is being set correctly by mapBCallBack function", func(t *testing.T) {

		resetTestApiConfig()
		assertStrings(testApiConfig.Next, originalUrl, t)
		assertStrings(testApiConfig.Prev, "", t)

		//make tie first call to the API, should return noPrevEntriesText
		mapBCallBack(buf, &testApiConfig)
		assertStrings(trimString(replaceNewLines(buf.String())), noPrevEntriesText, t)
		resetBuffer()

		//first call map, then mapB
		mapCallBack(buf, &testApiConfig)
		mapBCallBack(buf, &testApiConfig)
		assertStrings(testApiConfig.Next, offsetTwentyUrl, t)
		assertStrings(testApiConfig.Prev, "", t)
		resetBuffer()

	})
}

func resetTestApiConfig() {
	testApiConfig.Next = originalUrl
	testApiConfig.Prev = ""
}

func resetBuffer() {
	buf.Reset()
}

func assertValues(got, want CliCommand, t testing.TB) {
	t.Helper()
	isCommandDescMatch := getCommandString(got) == getCommandString(want)
	isFuncPointersMatch := getPointer(got.Callback) == getPointer(want.Callback)

	if !isCommandDescMatch || !isFuncPointersMatch {
		t.Errorf("\ngot: %v\nwant %v", got, want)
	}
}

func generateExpectedPrint(input []string) (returnedString string) {
	for _, text := range input {
		returnedString += text + "\n"
	}
	return returnedString
}

func getCommandString(commandIn CliCommand) string {
	returnedString := fmt.Sprintf("%s|%s", commandIn.Name, commandIn.Description)
	return returnedString
}

func getPointer(inputFunc func(io.Writer, *apiIntegration.Config) bool) (pointer uintptr) {
	return reflect.ValueOf(inputFunc).Pointer()
}

func assertStrings(got, want string, t testing.TB) {
	t.Helper()
	if got != want {
		t.Errorf("\ngot:\n%s\nwant:\n%s\n", got, want)
	}
}

func assertStringContains(got, want string, t testing.TB) {
	t.Helper()
	if !strings.Contains(got, want) {
		t.Errorf("%q could not be found in %q", got, want)
	}
}
