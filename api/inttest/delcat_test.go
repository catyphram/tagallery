package inttest

import (
	"context"
	"reflect"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"tagallery.com/api/config"
	"tagallery.com/api/model"
	"tagallery.com/api/mongodb"
	"tagallery.com/api/testutil"
)

var categoryFixture = model.Category{
	Name:        "Category 1",
	Description: "Category 1 description.",
}
var insertID string

func createCategoryDBFixturesForDeletion(ctx context.Context, db string) error {
	collection := mongodb.Client().Database(db).Collection("category")
	result, err := collection.InsertOne(ctx, categoryFixture, options.InsertOne())
	if err != nil {
		return err
	}
	insertID = result.InsertedID.(primitive.ObjectID).Hex()
	return nil
}

func deleteCategorySetup(t *testing.T) func() {
	if err := createCategoryDBFixturesForDeletion(context.Background(), config.Get().Database); err != nil {
		format, args := testutil.FormatTestError(
			"An error while creating the database fixtures has occured.",
			map[string]interface{}{
				"error": err,
			})
		t.Fatalf(format, args...)
	}
	return func() {
		testutil.CleanCollection(t, config.Get().Database, "category")
	}
}

func DeleteCategory(t *testing.T) {
	var categories []model.Category
	var response interface{}

	defer deleteCategorySetup(t)()

	if err := DeleteRequest(apiURL("/category/"+insertID), &response); err != nil {
		format, args := testutil.FormatTestError(
			"Request failed.",
			map[string]interface{}{
				"error":    err,
				"insertID": insertID,
			})
		t.Errorf(format, args...)
	}

	if err := GetRequest(apiURL("/category"), &categories); err != nil {
		format, args := testutil.FormatTestError(
			"Request failed.",
			map[string]interface{}{
				"error": err,
			})
		t.Errorf(format, args...)
	}

	if !reflect.DeepEqual([]model.Category{}, categories) {
		format, args := testutil.FormatTestError(
			"Deleted category should not be returned via request.",
			map[string]interface{}{
				"expected": []model.Category{},
				"got":      response,
			})
		t.Errorf(format, args...)
	}
}
