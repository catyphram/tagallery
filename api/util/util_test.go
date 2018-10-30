package util

import (
	"os"
	"path/filepath"
	"testing"

	"tagallery.com/api/testutil"
)

func TestJoin(t *testing.T) {
	expected := "Hello World"
	str := Join("Hello", " ", "World")
	if str != expected {
		format, args := testutil.FormatTestError(
			"TestJoin failed to join strings.",
			map[string]interface{}{
				"expected": expected,
				"got":      str,
			})
		t.Errorf(format, args...)
	}
}

// createTestDirectory creates and fills a directory with some example content.
func createTestDirectory(directory string) error {

	if err := os.Mkdir(directory, os.ModePerm); err != nil {
		return err
	}

	return testutil.FillDirectory(directory, testutil.FileTree{
		Dirs: map[string]*testutil.FileTree{
			"nested": &testutil.FileTree{},
		},
		Files: []string{"example.txt"},
	})
}

// readDirRecursive recursively reads the names of paths in a directory and returns them.
func readDirRecursive(dir string) ([]string, error) {
	var content []string
	walkErr := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err == nil {
			if path != dir {
				content = append(content, path)
			}
			return nil
		} else {
			return err
		}
	})
	return content, walkErr
}

func TestEmptyDirectory(t *testing.T) {
	var dir = "testdata"

	if err := createTestDirectory(dir); err != nil {
		format, args := testutil.FormatTestError(
			"Unable to create test directory for deletion.",
			map[string]interface{}{
				"error": err,
			})
		t.Errorf(format, args...)
	}

	if err := EmptyDirectory(dir); err != nil {
		format, args := testutil.FormatTestError(
			"Unable to empty directory.",
			map[string]interface{}{
				"error": err,
			})
		t.Errorf(format, args...)
	}

	if content, err := readDirRecursive(dir); err != nil || len(content) > 0 {
		format, args := testutil.FormatTestError(
			"Some content still resides in the directory.",
			map[string]interface{}{
				"error":   err,
				"content": content,
			})
		t.Errorf(format, args...)
	}

	os.RemoveAll(dir)

}
