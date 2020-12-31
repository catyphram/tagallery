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
	var response model.Category

	defer postCategorySetup(t)()

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
				"expected": category,
				"got":      response,
			})
		t.Errorf(format, args...)
	}
}
