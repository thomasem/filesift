package commands

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	KeepFirst  = "f"
	KeepSecond = "s"
	KeepBoth   = "b"
)

const (
	FindSimilarCommandName = "find-similar"
)

type FindSimilarCommand struct {
	flagSet *flag.FlagSet
	src     string
}

func NewFindSimilarCommand() *FindSimilarCommand {
	cmd := &FindSimilarCommand{
		flagSet: flag.NewFlagSet(FindSimilarCommandName, flag.ExitOnError),
	}

	cmd.flagSet.StringVar(&cmd.src, SourceDirectoryOptionName, DefaultSourceDir, "source directory")

	return cmd
}

func (dc *FindSimilarCommand) Run(args []string) error {
	if err := dc.flagSet.Parse(args); err != nil {
		return err
	}

	seen := []string{}

	err := filepath.WalkDir(dc.src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.Type().IsRegular() {
			return nil
		}

		if len(seen) > 1 {
			for i := range seen {
				// we should probably handle deletion here, so we can also remove from the seen list
				kd, err := compareFiles(seen[i], path)
				if err != nil {
					return err
				}
				fmt.Println("kd: ", kd)
			}
		}
		seen = append(seen, path)
		return nil
	})
	return err
}

func keepPrompt() string {
	var response string
	fmt.Print("Please specify which file(s) to keep. Enter (f)irst / (s)econd / (B)oth: ")
	fmt.Scanln(&response)

	switch strings.ToLower(response) {
	case KeepFirst, "first":
		return KeepFirst
	case KeepSecond, "second":
		return KeepSecond
	default:
		return KeepBoth
	}
}

func openToView(path string) error {
	switch runtime.GOOS {
	case "windows":
		return exec.Command("cmd", "/c", "start", path).Run() // UNTESTED
	case "darwin":
		return exec.Command("open", path).Run() // UNTESTED
	case "linux":
		return exec.Command("xdg-open", path).Run()
	default:
		return fmt.Errorf("unsupported operating system")
	}
}

func compareFiles(path1, path2 string) (string, error) {
	file1, err := os.Open(path1)
	if err != nil {
		return "", err
	}
	defer file1.Close()
	file2, err := os.Open(path2)
	if err != nil {
		return "", err
	}
	defer file2.Close()
	// For now, let's assume it's evaluating to TRUE for testing; we can implement comparison logic later
	// Print files to console
	fmt.Println("Found similar files:")
	fmt.Println(path1)
	fmt.Println(path2)

	if err := openToView(path1); err != nil {
		return "", err
	}
	if err := openToView(path2); err != nil {
		return "", err
	}

	return keepPrompt(), nil
}
