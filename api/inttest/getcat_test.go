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

var insertedCategories []model.Category
var categoryFixtures = []model.Category{
	{Name: "Category 1", Description: "Category 1 description."},
	{Name: "Category 2", Description: "Category 2 description."},
	{Name: "Category 3", Description: "Category 3 description."},
	{Name: "Category 4", Description: "Category 4 description."},
	{Name: "Category 5", Description: "Category 5 description."},
	{Name: "Category 6", Description: "Category 6 description."},
}

func createCategoryDBFixtures(ctx context.Context, db string) error {
	collection := mongodb.Client().Database(db).Collection("category")
	categories := make([]interface{}, len(categoryFixtures))
	for k, v := range categoryFixtures {
		categories[k] = v
	}

	result, err := collection.InsertMany(ctx, categories, options.InsertMany())
	if err != nil {
		return err
	}

	for i, id := range result.InsertedIDs {
		insertID := id.(primitive.ObjectID).Hex()
		insertedCategories = append(insertedCategories, model.Category{
			ID:          &insertID,
			Name:        categoryFixtures[i].Name,
			Description: categoryFixtures[i].Description,
		})
	}
	return nil
}

func getCategoriesSetup(t *testing.T) func() {
	if err := createCategoryDBFixtures(context.Background(), config.Get().Database); err != nil {
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

func GetCategories(t *testing.T) {
	var categories []model.Category

	defer getCategoriesSetup(t)()

	if err := GetRequest(apiURL("/category"), &categories); err != nil {
		format, args := testutil.FormatTestError(
			"Request failed.",
			map[string]interface{}{
				"error": err,
			})
		t.Errorf(format, args...)
	}

	if !reflect.DeepEqual(categories, insertedCategories) {
		format, args := testutil.FormatTestError(
			"Returned categories do not match expectations.",
			map[string]interface{}{
				"expected": categoryFixtures,
				"got":      categories,
			})
		t.Errorf(format, args...)
	}
}
