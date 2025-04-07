package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/thomasem/filesift/internal/copy"
	"github.com/thomasem/filesift/internal/unique"
)

const (
	DEFAULT_OUTPUT_DIR = "output"
	DEFAULT_SOURCE_DIR = "."

	FILE_MODE = 0755
)

var (
	ErrAborted = errors.New("operation aborted")
)

func confirm() bool {
	var response string
	fmt.Print("Enter [y/N]: ")
	fmt.Scanln(&response)
	return response == "y" || response == "Y"
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
		return false, fmt.Errorf("output path (%s) is not a directory", path)
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

		err = os.RemoveAll(dst)
		if err != nil {
			return err
		}
	}

	return os.Mkdir(dst, FILE_MODE)
}

// Few things we need here:
// 1. Ensure that the copied files have the same permissions as the original files
// 2. Ensure that the copied files have the same ownership as the original files
// 3. Ensure that the copied files have the same timestamps as the original files - NOT DOING
// 4. Ensure that copied files have same path as original (appreciating that duplicates will be ignored after the first
// copy is found)
// 5. Allow for in-place, which would introduce roughly the same way for finding unique files, but we would delete any
// duplicates found in the source directory instead of copying to an output directory. Right now, I prefer copying to
// preserve the source files
// 6. Support an ignore file (such as for ignoring specific types of files), which resolves up the filesystem tree, like .gitignore
// 7. Make sure to identify that this evaluates the contents of files, not the filenames. That is an edge case I'd need to handle, though - when storing, what happens with filename collisions?
// 8. Consider goroutine pool for parallel copying
func main() {
	src := flag.String("dir", DEFAULT_SOURCE_DIR, "source directory")
	dst := flag.String("out", DEFAULT_OUTPUT_DIR, "output directory")
	dryRun := flag.Bool("dry-run", false, "dry run")
	flag.Parse()

	uniqueFiles, err := unique.GetUniqueFiles(*src)
	if err != nil {
		fmt.Println("error evaluating files:", err)
		os.Exit(1)
	}

	if *dryRun {
		for _, path := range uniqueFiles {
			fmt.Println(path)
		}
		return
	}

	err = ensureOutputDirectory(*dst)
	if errors.Is(err, ErrAborted) {
		fmt.Println("Operation aborted.")
		os.Exit(0)
	}
	if err != nil {
		fmt.Println("error setting up output directory:", err)
		os.Exit(1)
	}

	for _, path := range uniqueFiles {
		dstPath := filepath.Join(*dst, filepath.Base(path))
		if err := copy.CopyFile(path, dstPath); err != nil {
			fmt.Printf("error copying file %s to %s: %v\n", path, dstPath, err)
			os.Exit(1)
		}
	}
}
