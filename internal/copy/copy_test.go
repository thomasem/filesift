package copy

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/thomasem/filesift/internal/testdata"
	"github.com/thomasem/filesift/internal/unique"
)

func TestCopyFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "test-copy-file-*")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	testFile := filepath.Join(testdata.TestDataDir, "elsa_sploot.jpg")
	destFile := filepath.Join(tmpDir, "elsa_sploot.jpg")

	testFileInfo, err := os.Stat(testFile)
	assert.NoError(t, err)

	err = CopyFile(testFile, destFile)
	assert.NoError(t, err)

	copyFileInfo, err := os.Stat(destFile)
	assert.NoError(t, err)

	assert.Equal(t, testFileInfo.Name(), copyFileInfo.Name())
	assert.Equal(t, testFileInfo.Size(), copyFileInfo.Size())
	assert.Equal(t, testFileInfo.Mode(), copyFileInfo.Mode())

	tfSum, err := unique.ComputeChecksum(testFile)
	assert.NoError(t, err)
	cfSum, err := unique.ComputeChecksum(destFile)
	assert.NoError(t, err)
	assert.Equal(t, tfSum, cfSum)
}
