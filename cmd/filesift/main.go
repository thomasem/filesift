package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/thomasem/filesift/internal/commands"
)

func showHelp() {
	fmt.Printf(`Usage: filesift <command> [options]

Commands:
    dedupe    Find duplicate or similar files

Use "filesift <command> -h" for help with a specific command.
`)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Missing required argument(s)")
		showHelp()
		os.Exit(1)
	}

	firstArg, rest := os.Args[1], os.Args[2:]
	if firstArg == "-h" || firstArg == "--help" {
		showHelp()
		os.Exit(0)
	}

	var err error

	switch firstArg {
	case commands.DeduplicateCommandName:
		cmd := commands.NewDeduplicateCommand()
		err = cmd.Run(rest)
	default:
		fmt.Println("Command not found")
		showHelp()
		os.Exit(1)
	}

	if err != nil {
		fmt.Println(err)
		if errors.Is(err, commands.ErrAborted) {
			os.Exit(0)
		}
		os.Exit(1)
	}
}
