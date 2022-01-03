package horcrux

import (
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestGetSplitFilename(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		output   string
		id       int
		expected string
	}{
		{
			"default path",
			"/a/b/c.md",
			"",
			1,
			"/a/b/c",
		},
		{
			"default path",
			"/a/b/c.md",
			"/d",
			1,
			"/d/c",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := getSplitFilebase(tt.filename, tt.output)
			if actual != tt.expected {
				t.Errorf("getSplitFilename() = %v, expected %v", actual, tt.expected)
			}
		})
	}
}

func TestSplit(t *testing.T) {
	tempDir := t.TempDir()

	// Write a plaintext file
	plaintext := "secret plaintext"
	plaintextFilename := path.Join(tempDir, "plaintext")
	err := os.WriteFile(plaintextFilename, []byte(plaintext), 0644)
	Assert(err)

	// Split it
	shareDir := t.TempDir()
	err = Split(plaintextFilename, 3, 2, shareDir)
	Assert(err)

	// Restore it,
	files, err := ioutil.ReadDir(shareDir)
	fileNames := []string{}
	for _, f := range files {
		fileNames = append(fileNames, path.Join(shareDir, f.Name()))
	}
	Assert(err)
	original, err := Restore(fileNames)
	Assert(err)

	if string(original) != plaintext {
		t.Errorf("Split(): %b != %s", original, plaintext)
	}
}
