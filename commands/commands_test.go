package commands

import (
	"fmt"
	"reflect"
	"testing"
)

type getCLICommandReturnVal struct {
	command CliCommand
	err     error
}

type testStruct struct {
	name     string
	input    string
	expected getCLICommandReturnVal
}

func TestGetCLICommand(t *testing.T) {
	tests := []testStruct{
		{
			name:  "testing help command",
			input: "help",
			expected: getCLICommandReturnVal{
				helpCommand, nil,
			},
		},
		{
			name:  "testing no command",
			input: "erroneous text",
			expected: getCLICommandReturnVal{
				CliCommand{}, fmt.Errorf("key %s is not valid", "erroneous text"),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			command, err := GetCLICommand(test.input)
			got := getCLICommandReturnVal{
				command, err,
			}

			assertValues(got, test.expected, t)

		})
	}
}

func assertValues(got, want getCLICommandReturnVal, t testing.TB) {
	t.Helper()
	isCommandDescMatch := getCommandString(got.command) == getCommandString(want.command)
	isFuncPointersMatch := getPointer(got.command.Callback) == getPointer(want.command.Callback)
	isErrorsMatch := compareErrorStrings(got.err, want.err)
	if !isCommandDescMatch || !isFuncPointersMatch || !isErrorsMatch {
		t.Errorf("\ngot: %v\nwant %v", got, want)
	}
}

func getCommandString(commandIn CliCommand) string {
	returnedString := fmt.Sprintf("%s|%s", commandIn.Name, commandIn.Description)
	return returnedString
}

func getPointer(inputFunc func() error) (pointer uintptr) {
	return reflect.ValueOf(inputFunc).Pointer()
}

func compareErrorStrings(gotError, wantError error) bool {
	gotString := getErrorString(gotError)
	wantString := getErrorString(wantError)
	return gotString == wantString
}

func getErrorString(err error) (errorString string) {
	if err == nil {
		return ""
	}
	return err.Error()
}
