package testutil

import (
	"context"
	"testing"

	"tagallery.com/api/mongodb"
)

// DropCollection deletes all collections in a MongoDB database.
func DropCollection(db string, collection string) error {
	client := mongodb.Client()
	return client.Database(db).Collection(collection).Drop(context.Background())
}

// CleanCollection drops a db via DropCollection() and then logs potential errors.
func CleanCollection(t *testing.T, db string, collection string) {
	if err := DropCollection(db, collection); err != nil {
		format, args := FormatTestError(
			"Failed to drop the database.",
			map[string]interface{}{
				"error": err,
			})
		t.Errorf(format, args...)
	}
}
