package unique

import (
	"path/filepath"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thomasem/filesift/internal/testdata"
)

func TestComputeChecksum(t *testing.T) {
	tests := []struct {
		file string
		sum  string
	}{
		{"text_a", "c5089fa8d41e8e3990ee37e7cf673e992afcee7c2b6ff9c3498a7352ed81fc7e"},
		{"text_a_duplicate", "c5089fa8d41e8e3990ee37e7cf673e992afcee7c2b6ff9c3498a7352ed81fc7e"},
		{"text_b", "56cb6e3f37d0f05499e4ceee318b97a49931520430cfa9b33cffe52b1b11311d"},
		{"elsa_sploot.jpg", "0a5536af7d19d7d7efcb2e0a642416a483f93176b5c89fac384b850a61ef2ff8"},
		{"elsa_sploot_duplicate.jpg", "0a5536af7d19d7d7efcb2e0a642416a483f93176b5c89fac384b850a61ef2ff8"},
		{"elsa_walk.jpg", "f0ed2952bf1935b5415fa69e5451a74956899b0375e9e1a70e0d30cd3b4fc49d"},
		{"bumble.png", "98d37deddecdbaa222e11798f9c1e96a02f207bf6e40132b0450153576de2ae5"},
		{"bumble_duplicate.png", "98d37deddecdbaa222e11798f9c1e96a02f207bf6e40132b0450153576de2ae5"},
		{"flower.jpg", "33fbd81de3097c9fa377a6af63202699bc71060e525f89ca69cac26cfa0d0169"},
	}

	for _, test := range tests {
		t.Run(test.file, func(t *testing.T) {
			sum, err := ComputeChecksum(filepath.Join(testdata.TestDataDir, test.file))
			assert.NoError(t, err)
			assert.Equal(t, sum, test.sum)
		})
	}
}

func TestGetUniqueFiles(t *testing.T) {
	expectedUniqueFiles := []string{
		"text_a",
		"text_b",
		"elsa_sploot.jpg",
		"elsa_walk.jpg",
		"bumble.png",
		"flower.jpg",
	}

	for i, f := range expectedUniqueFiles {
		expectedUniqueFiles[i] = filepath.Join(testdata.TestDataDir, f)
	}

	uniqueFiles, err := GetUniqueFiles(testdata.TestDataDir)
	assert.NoError(t, err)
	sort.Strings(expectedUniqueFiles)
	sort.Strings(uniqueFiles)
	assert.Equal(t, uniqueFiles, expectedUniqueFiles)
}
