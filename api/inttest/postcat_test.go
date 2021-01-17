package inttest

import (
	"reflect"
	"testing"

	"tagallery.com/api/config"
	"tagallery.com/api/model"
	"tagallery.com/api/testutil"
)

func postCategorySetup(t *testing.T) func() {
	return func() {
		testutil.CleanCollection(t, config.Get().Database, "category")
	}
}

func PostCategory(t *testing.T) {
	var categories []model.Category
	var category = model.Category{
		Name:        "Category 1",
		Description: "Test Category",
	}
	var insertID string
	var response model.Category
	var errorResponse ErrorResponse

	defer postCategorySetup(t)()

	if err := PostRequest(apiURL("/category"), category, &response); err != nil {
		format, args := testutil.FormatTestError(
			"Request failed.",
			map[string]interface{}{
				"error": err,
			})
		t.Errorf(format, args...)
	}

	insertID = *response.ID
	category.ID = &insertID

	if !reflect.DeepEqual(category, response) {
		format, args := testutil.FormatTestError(
			"Returned category does not match inserted one.",
			map[string]interface{}{
				"expected": category,
				"got":      response,
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

	if !reflect.DeepEqual([]model.Category{category}, categories) {
		format, args := testutil.FormatTestError(
			"Inserted category is not returned via request.",
			map[string]interface{}{
				"expected": []model.Category{category},
				"got":      response,
			})
		t.Errorf(format, args...)
	}

	category.Name = "Category 2"
	category.Description = "Category 2 Description"

	if err := PostRequest(apiURL("/category"), category, &response); err != nil {
		format, args := testutil.FormatTestError(
			"Request failed.",
			map[string]interface{}{
				"error": err,
			})
		t.Errorf(format, args...)
	}

	if !reflect.DeepEqual(category, response) {
		format, args := testutil.FormatTestError(
			"Returned category does not match updated one.",
			map[string]interface{}{
				"expected": category,
				"got":      response,
			})
		t.Errorf(format, args...)
	}

	category.ID = nil
	category.Description = "Updated description"

	if err := PostRequest(apiURL("/category"), category, &response); err != nil {
		format, args := testutil.FormatTestError(
			"Request failed.",
			map[string]interface{}{
				"error": err,
			})
		t.Errorf(format, args...)
	}

	category.ID = &insertID

	if !reflect.DeepEqual(category, response) {
		format, args := testutil.FormatTestError(
			"Returned category does not match updated one when updating using the name as identifier.",
			map[string]interface{}{
				"expected": category,
				"got":      response,
			})
		t.Errorf(format, args...)
	}

	outdatedObjectID := "600393d56c57d714f7f1fe8f"
	category.ID = &outdatedObjectID

	if err := PostRequest(apiURL("/category"), category, &errorResponse); err != nil ||
		errorResponse.Error == "" {
		format, args := testutil.FormatTestError(
			"Request failed or no error was returned for a non-unique name insert.",
			map[string]interface{}{
				"error": err,
			})
		t.Errorf(format, args...)
	}

	invalidObjectID := "abc123"
	category.ID = &invalidObjectID

	if err := PostRequest(apiURL("/category"), category, &errorResponse); err != nil ||
		errorResponse.Error == "" {
		format, args := testutil.FormatTestError(
			"Request failed or no error was returned for an invalid id insert.",
			map[string]interface{}{
				"error": err,
			})
		t.Errorf(format, args...)
	}
}
