package unique

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func ComputeChecksum(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	sum := sha256.New()
	if _, err := io.Copy(sum, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(sum.Sum(nil)), nil
}

func GetUniqueFiles(dir string) ([]string, error) {
	uniqueFiles := make(map[string]string)
	// WalkDir is more efficient by avoiding extra system calls (https://pkg.go.dev/path/filepath#WalkDir)
	err := filepath.WalkDir(dir, func(path string, info os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// info.IsDir() isn't enough because there are several other types of filesystem entries
		// that aren't regular files and trying to calculate a hash for them will fail.
		if !info.Type().IsRegular() {
			return nil
		}

		sum, err := ComputeChecksum(path)
		if err != nil {
			return fmt.Errorf("error computing checksum for %s: %s", path, err)
		}
		if _, ok := uniqueFiles[sum]; !ok {
			uniqueFiles[sum] = path
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	paths := make([]string, 0, len(uniqueFiles))
	for _, file := range uniqueFiles {
		paths = append(paths, file)
	}

	return paths, nil
}
