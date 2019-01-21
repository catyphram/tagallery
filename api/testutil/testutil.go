// Package testutil contains util & helper functions used by tests.
package testutil

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mongodb/mongo-go-driver/mongo"
)

// Join concatenates multiple strings into one.
// Same as in the util package, but has to be duplicated due to cycle import otherwise.
func join(strs ...string) string {
	var sb strings.Builder
	for _, str := range strs {
		sb.WriteString(str)
	}
	return sb.String()
}

// FormatTestError creates a formatted string to pass to log functions.
// The function returns a interpolated string and the values for this string.
func FormatTestError(desc string, args map[string]interface{}) (string, []interface{}) {
	var params []interface{}
	logString := []string{desc}

	for k, v := range args {
		params = append(params, v)
		logString = append(logString, "\n", k, ": %v")
	}

	return join(logString...), params
}

// FileTree represents a level of a directory structure.
type FileTree struct {
	Dirs  map[string]*FileTree
	Files []string
}

// TouchFile creates an empty file with permission 0666.
func TouchFile(path string) error {
	_, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0666)
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
		if err := os.Mkdir(subDir, os.ModePerm); err != nil {
			return err
		}
		if err := FillDirectory(subDir, *content); err != nil {
			return err
		}
	}

	return nil
}

// OpenMongoDbConnection returns a connection to a mongo db.
func OpenMongoDbConnection(host string, database string) (*mongo.Database, error) {
	uri := fmt.Sprintf(`mongodb://%s`, host)

	client, err := mongo.NewClient(uri)
	if err != nil {
		return nil, err
	}

	err = client.Connect(nil)
	if err != nil {
		return nil, err
	}

	db := client.Database(database)
	return db, nil
}

// CleanMongoDbFromConnection drops all collections in mongo db.
// Same as CleanMongoDb, but requires *mongo.Database as an argument.
func CleanMongoDbFromConnection(db *mongo.Database) error {
	if db == nil {
		return errors.New("db is nil")
	}
	if err := db.Drop(nil); err != nil {
		return err
	}
	return nil
}

// CleanMongoDb drops all collections in mongo db.
// Same as CleanMongoDbFromConnection, but requires host and database as arugments.
func CleanMongoDb(host string, database string) error {
	if db, err := OpenMongoDbConnection(host, database); err != nil {
		return err
	} else {
		if err := CleanMongoDbFromConnection(db); err != nil {
			return err
		}
	}
	return nil
}

// DropMongoDb calls CleanMongoDb and prints a formatted error message
func DropMongoDb(host string, database string, t *testing.T) {
	err := CleanMongoDb(host, database)
	if err != nil {
		format, args := FormatTestError(
			"Failed to clean up the database after testing.",
			map[string]interface{}{
				"error": err,
			})
		t.Errorf(format, args...)
	}
}

// InsertIntoMongoDb connects to the mongodb and inserts the passed content.
func InsertIntoMongoDb(host string, database string, collection string, content []interface{}) error {
	db, err := OpenMongoDbConnection(host, database)

	if err != nil {
		return err
	}
	defer db.Client().Disconnect(context.Background())

	if err != nil {
		return err
	}

	_, err = db.Collection(collection).InsertMany(nil, content)
	if err != nil {
		return err
	}

	return nil
}
