package commands

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"testing"
)

type testStruct struct {
	name                   string
	input                  string
	expectedStringsInprint []string
	expected               CliCommand
}

func TestGetCLICommand(t *testing.T) {
	tests := []testStruct{
		{
			name:                   "testing help command",
			input:                  "help",
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
		{
			name:                   "testing map command",
			input:                  "map",
			expected:               mapCommand,
			expectedStringsInprint: []string{mapString},
		},
		{
			name:                   "testing mapB command",
			input:                  "mapb",
			expected:               mapBCommand,
			expectedStringsInprint: []string{mapBString},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := GetCLICommand(test.input)
			assertValues(got, test.expected, t)
			buf := &bytes.Buffer{}
			got.Callback(buf)
			assertStrings(buf.String(), generateExpectedPrint(test.expectedStringsInprint), t)
		})
	}
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

func getPointer(inputFunc func(io.Writer) bool) (pointer uintptr) {
	return reflect.ValueOf(inputFunc).Pointer()
}

func assertStrings(got, want string, t testing.TB) {
	t.Helper()
	if got != want {
		t.Errorf("\ngot:\n%s\nwant:\n%s\n", got, want)
	}
}
