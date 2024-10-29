package commands

import "fmt"

type CliCommand struct {
	Name        string
	Description string
	Callback    func() error
}

var (
	helpCommand CliCommand = CliCommand{
		"help",
		"this is a command to describe help functions",
		func() error {
			return nil
		},
	}
	exitCommand CliCommand = CliCommand{
		"exit",
		"this is to exit the program",
		func() error {
			return nil
		},
	}
)

var CommandMap = map[string]CliCommand{
	"help": helpCommand,
	"exit": exitCommand,
}

func GetCLICommand(text string) (command CliCommand, err error) {
	command, ok := CommandMap[text]
	if !ok {
		return CliCommand{}, fmt.Errorf("key %s is not valid", text)
	}
	return command, nil
}
