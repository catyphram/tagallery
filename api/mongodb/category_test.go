package mongodb_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"tagallery.com/api/config"
	"tagallery.com/api/logger"
	"tagallery.com/api/model"
	"tagallery.com/api/mongodb"
	"tagallery.com/api/testutil"
)

var testCategory = model.Category{
	Name:        "Test Category",
	Description: "This is a test category.",
}

func init() {
	logger.Setup(true)
}

func TestUpsertCategory(t *testing.T) {
	configuration := config.Load()

	defer testutil.CleanCollection(t, configuration.Database, "category")
	mongodb.Connect(context.Background(), fmt.Sprintf(`mongodb://%s`, config.Get().DatabaseHost))

	err := mongodb.UpsertCategory(testCategory)

	if err != nil {
		format, args := testutil.FormatTestError(
			"Expected category to be inserted.",
			map[string]interface{}{
				"error": err,
			})
		t.Errorf(format, args...)
	}
}

func TestQueryCategories(t *testing.T) {
	configuration := config.Load()

	defer testutil.CleanCollection(t, configuration.Database, "category")
	mongodb.Connect(context.Background(), fmt.Sprintf(`mongodb://%s`, config.Get().DatabaseHost))

	if err := mongodb.UpsertCategory(testCategory); err != nil {
		format, args := testutil.FormatTestError(
			"Failed to create test category.",
			map[string]interface{}{
				"error": err,
			})
		t.Errorf(format, args...)
	}

	categories, err := mongodb.QueryCategories()

	if err != nil ||
		!reflect.DeepEqual(categories, []model.Category{testCategory}) {
		format, args := testutil.FormatTestError(
			"Expected category to be inserted.",
			map[string]interface{}{
				"error": err,
			})
		t.Errorf(format, args...)
	}
}

func TestDeleteCategory(t *testing.T) {
	configuration := config.Load()

	defer testutil.CleanCollection(t, configuration.Database, "category")
	mongodb.Connect(context.Background(), fmt.Sprintf(`mongodb://%s`, config.Get().DatabaseHost))

	collection := mongodb.Client().Database(configuration.Database).Collection("category")
	result, err := collection.InsertOne(context.Background(), testCategory, options.InsertOne())
	insertId := result.InsertedID.(primitive.ObjectID).Hex()

	if err != nil {
		format, args := testutil.FormatTestError(
			"Failed to create test category.",
			map[string]interface{}{
				"error": err,
			})
		t.Errorf(format, args...)
	}

	if err := mongodb.DeleteCategory(insertId); err != nil {
		format, args := testutil.FormatTestError(
			"Expected category to be deleted.",
			map[string]interface{}{
				"error": err,
			})
		t.Errorf(format, args...)
	}

	if result := collection.FindOne(context.Background(), bson.M{"_id": insertId}); result.Err() != mongo.ErrNoDocuments {
		format, args := testutil.FormatTestError(
			"Expected category to be removed from the database.",
			map[string]interface{}{
				"error": result.Err(),
			})
		t.Errorf(format, args...)
	}
}
