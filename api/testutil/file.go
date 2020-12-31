// Package testutil contains util & helper functions used by tests.
package testutil

import (
	"os"
	"path/filepath"
)

// FileTree represents a level of a directory structure.
type FileTree struct {
	Dirs  map[string]*FileTree
	Files []string
}

// TouchFile creates an empty file with permission 0755.
func TouchFile(path string) error {
	_, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0755)
	return err
}

// FillDirectory creates folders and touches files according to a given fileTree structure.
// The function calls itself recursively per directory.
func FillDirectory(directory string, fileTree FileTree) error {
	for _, file := range fileTree.Files {
		if err := TouchFile(filepath.Join(directory, file)); err != nil {
			return err
		}
	}

	for dir, content := range fileTree.Dirs {
		subDir := filepath.Join(directory, dir)
		if err := os.Mkdir(subDir, 0755); err != nil {
			return err
		}
		if err := FillDirectory(subDir, *content); err != nil {
			return err
		}
	}

	return nil
}
