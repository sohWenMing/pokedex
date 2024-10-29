package commands

import (
	"fmt"
	"reflect"
	"testing"
)

type testStruct struct {
	name     string
	input    string
	expected CliCommand
}

func TestGetCLICommand(t *testing.T) {
	tests := []testStruct{
		{
			name:     "testing help command",
			input:    "help",
			expected: helpCommand,
		},
		{
			name:     "testing no command",
			input:    "erroneous text",
			expected: defaultUsageCommand,
		},
		{
			name:     "testing exit command",
			input:    "exit",
			expected: exitCommand,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := GetCLICommand(test.input)
			assertValues(got, test.expected, t)

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

func getCommandString(commandIn CliCommand) string {
	returnedString := fmt.Sprintf("%s|%s", commandIn.Name, commandIn.Description)
	return returnedString
}

func getPointer(inputFunc func() bool) (pointer uintptr) {
	return reflect.ValueOf(inputFunc).Pointer()
}
