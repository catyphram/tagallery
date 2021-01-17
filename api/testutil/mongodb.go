package testutil

import (
	"context"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"tagallery.com/api/mongodb"
)

// DropCollection deletes all collections in a MongoDB database.
func DropCollection(db string, collection string) error {
	client := mongodb.Client()
	_, err := client.Database(db).Collection(collection).DeleteMany(context.Background(), bson.M{}, options.Delete())
	return err
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
