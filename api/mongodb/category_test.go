package mongodb_test

import (
	"context"
	"errors"
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

	category, err := mongodb.UpsertCategory(testCategory)

	if err != nil ||
		category.Name != testCategory.Name ||
		category.Description != testCategory.Description {
		format, args := testutil.FormatTestError(
			"Expected category to be inserted.",
			map[string]interface{}{
				"error": err,
			})
		t.Errorf(format, args...)
	}

	newCategory := model.Category{
		ID:          category.ID,
		Name:        "New category",
		Description: "New category description",
	}

	updatedCategory, err := mongodb.UpsertCategory(newCategory)
	if err != nil ||
		updatedCategory.Name != newCategory.Name ||
		updatedCategory.Description != newCategory.Description {
		format, args := testutil.FormatTestError(
			"Expected same category to be updated.",
			map[string]interface{}{
				"error":    err,
				"expected": newCategory,
				"got":      updatedCategory,
			})
		t.Errorf(format, args...)
	}

	invalidObjectID := "123xzy"
	invalidCategory := model.Category{
		ID:          &invalidObjectID,
		Name:        "Invalid category",
		Description: "Invalid category description",
	}
	_, err = mongodb.UpsertCategory(invalidCategory)
	if err == nil {
		format, args := testutil.FormatTestError(
			"Expected category with an invalid object id to not be updated.",
			map[string]interface{}{
				"error":    err,
				"category": invalidCategory,
			})
		t.Errorf(format, args...)
	}
}

func TestQueryCategories(t *testing.T) {
	configuration := config.Load()

	defer testutil.CleanCollection(t, configuration.Database, "category")
	mongodb.Connect(context.Background(), fmt.Sprintf(`mongodb://%s`, config.Get().DatabaseHost))

	category, err := mongodb.UpsertCategory(testCategory)
	if err != nil {
		format, args := testutil.FormatTestError(
			"Failed to create test category.",
			map[string]interface{}{
				"error": err,
			})
		t.Errorf(format, args...)
	}

	categories, err := mongodb.QueryCategories()
	if err != nil ||
		!reflect.DeepEqual(categories, []model.Category{*category}) {
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

	if err != nil {
		format, args := testutil.FormatTestError(
			"Failed to create test category.",
			map[string]interface{}{
				"error": err,
			})
		t.Errorf(format, args...)
	}
	insertId := result.InsertedID.(primitive.ObjectID).Hex()

	if err := mongodb.DeleteCategory(insertId); err != nil {
		format, args := testutil.FormatTestError(
			"Expected category to be deleted.",
			map[string]interface{}{
				"error": err,
			})
		t.Errorf(format, args...)
	}

	if result := collection.FindOne(context.Background(), bson.M{"_id": insertId}); !errors.Is(result.Err(), mongo.ErrNoDocuments) {
		format, args := testutil.FormatTestError(
			"Expected category to be removed from the database.",
			map[string]interface{}{
				"error": result.Err(),
			})
		t.Errorf(format, args...)
	}

	invalidInsertId := "123xyz"

	if err := mongodb.DeleteCategory(invalidInsertId); err == nil {
		format, args := testutil.FormatTestError(
			"Expected providing an invalid object id to fail.",
			map[string]interface{}{
				"error": err,
			})
		t.Errorf(format, args...)
	}
}
