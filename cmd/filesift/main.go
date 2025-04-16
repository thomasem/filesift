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
    dedupe    		Find duplicate or similar files
    find-similar 	Find similar files and prompt for which to keep

Use "filesift <command> -h" for help with a specific command.
`)
}

func exitWithHelp(msg string, code int) {
	if msg != "" {
		fmt.Println(msg)
	}
	showHelp()
	os.Exit(code)
}

func main() {
	if len(os.Args) < 2 {
		exitWithHelp("Missing required argument(s)", 1)
	}

	firstArg, rest := os.Args[1], os.Args[2:]
	if firstArg == "-h" || firstArg == "--help" {
		exitWithHelp("", 0)
	}

	var err error

	switch firstArg {
	case commands.DeduplicateCommandName:
		cmd := commands.NewDeduplicateCommand()
		err = cmd.Run(rest)
	case commands.FindSimilarCommandName:
		cmd := commands.NewFindSimilarCommand()
		err = cmd.Run(rest)
	default:
		exitWithHelp("Command not found", 1)
	}

	if err != nil {
		fmt.Println(err)
		if errors.Is(err, commands.ErrAborted) {
			os.Exit(0)
		}
		os.Exit(1)
	}
}
