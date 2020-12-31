package inttest

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"tagallery.com/api/config"
	"tagallery.com/api/model"
	"tagallery.com/api/testutil"
	"tagallery.com/api/util"
)

func postImageSetup(t *testing.T) func() {
	_ = os.MkdirAll(config.Get().Images, 0755)
	if err := createFileFixtures(
		config.Get().Images,
		[]string{"test.jpg"},
		[]string{},
	); err != nil {
		format, args := testutil.FormatTestError(
			"An error while creating the file fixtures has occured.",
			map[string]interface{}{
				"error": err,
			})
		t.Fatalf(format, args...)
	}

	return func() {
		testutil.CleanCollection(t, config.Get().Database, "image")
		os.RemoveAll(config.Get().Images)
	}
}

func PostImage(t *testing.T) {
	var images []model.Image
	var image = model.Image{
		File:               filepath.Join(config.Get().UnprocessedImagesFolder, "test.jpg"),
		AssignedCategories: []string{},
		ProposedCategories: []string{"Category 1"},
		StarredCategory:    util.StringPtr("Category 2"),
	}
	var updatedImage = model.Image{
		File:               filepath.Join(config.Get().ProcessedImagesFolder, "test.jpg"),
		AssignedCategories: []string{"Category 2"},
		ProposedCategories: []string{"Category 1"},
		StarredCategory:    util.StringPtr("Category 2"),
	}
	var response model.Image

	defer postImageSetup(t)()

	if err := PostRequest(apiURL("/image"), image, &response); err != nil {
		format, args := testutil.FormatTestError(
			"Request failed.",
			map[string]interface{}{
				"error": err,
			})
		t.Errorf(format, args...)
	}

	if !reflect.DeepEqual(updatedImage, response) {
		format, args := testutil.FormatTestError(
			"Returned image does not match inserted one.",
			map[string]interface{}{
				"expected": image,
				"got":      response,
			})
		t.Errorf(format, args...)
	}

	if err := GetRequest(apiURL("/image"), &images); err != nil {
		format, args := testutil.FormatTestError(
			"Request failed.",
			map[string]interface{}{
				"error": err,
			})
		t.Errorf(format, args...)
	}

	if !reflect.DeepEqual([]model.Image{updatedImage}, images) {
		format, args := testutil.FormatTestError(
			"Inserted image is not returned via request.",
			map[string]interface{}{
				"expected": image,
				"got":      response,
			})
		t.Errorf(format, args...)
	}
}
