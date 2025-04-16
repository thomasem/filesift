package commands

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/thomasem/filesift/internal/copy"
	"github.com/thomasem/filesift/internal/unique"
)

const (
	DeduplicateCommandName = "dedupe"

	SourceDirectoryOptionName = "src"
	OutputDirectoryOptionName = "out"
	DryRunOptionName          = "dry-run"
)

type DeduplicateCommand struct {
	flagSet *flag.FlagSet
	out     string
	src     string
	dryRun  bool
}

func NewDeduplicateCommand() *DeduplicateCommand {
	cmd := &DeduplicateCommand{
		flagSet: flag.NewFlagSet(DeduplicateCommandName, flag.ExitOnError),
	}

	cmd.flagSet.StringVar(&cmd.src, SourceDirectoryOptionName, DefaultSourceDir, "source directory")
	cmd.flagSet.StringVar(&cmd.out, OutputDirectoryOptionName, DefaultOutputDir, "output directory")
	cmd.flagSet.BoolVar(&cmd.dryRun, DryRunOptionName, false, "dry run")

	return cmd
}

func (dc *DeduplicateCommand) Run(args []string) error {
	if err := dc.flagSet.Parse(args); err != nil {
		return err
	}

	uniqueFiles, err := unique.GetUniqueFiles(dc.src)
	if err != nil {
		return fmt.Errorf("error evaluating files: %w", err)
	}

	if dc.dryRun {
		for _, path := range uniqueFiles {
			fmt.Println(path)
		}
		return nil
	}

	if err := ensureOutputDirectory(dc.out); err != nil {
		return err
	}

	for _, path := range uniqueFiles {
		dstPath := filepath.Join(dc.out, filepath.Base(path))
		if err := copy.CopyFile(path, dstPath); err != nil {
			return fmt.Errorf("error copying file %s to %s: %w", path, dstPath, err)
		}
	}
	return nil
}

func directoryExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if !info.IsDir() {
		return false, fmt.Errorf("path (%s) is not a directory", path)
	}
	return true, nil
}

func ensureOutputDirectory(dst string) error {
	exists, err := directoryExists(dst)
	if err != nil {
		return err
	}

	if exists {
		fmt.Printf("Output directory (%s) exists. Do you want to overwrite? ", dst)
		if !confirm() {
			return ErrAborted
		}

		if err = os.RemoveAll(dst); err != nil {
			return err
		}
	}

	return os.MkdirAll(dst, DefaultFileMode)
}
