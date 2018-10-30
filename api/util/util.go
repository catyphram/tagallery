// Package util contains utility functions that can be reused throughout the app.
package util

import (
	"os"
	"path/filepath"
	"strings"
)

// Join concatenates multiple strings into one.
func Join(strs ...string) string {
	var sb strings.Builder
	for _, str := range strs {
		sb.WriteString(str)
	}
	return sb.String()
}

// EmptyDirectory cleans the directory and leaves it empty. The root directory itself is kept.
func EmptyDirectory(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}
