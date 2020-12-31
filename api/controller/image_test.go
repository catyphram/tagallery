package controller_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"tagallery.com/api/config"
	"tagallery.com/api/controller"
	"tagallery.com/api/logger"
	"tagallery.com/api/model"
	"tagallery.com/api/mongodb"
	"tagallery.com/api/testutil"
	"tagallery.com/api/util"
)

// fileFixtures are the files created for the tests. The slice has to be sorted alphabetically
// otherwise the comparison will fail due to the different order.
var fileFixtures = []string{
	"1a8LsNsvVR.jpg",
	"8kGrdTdYgT.png",
	"gtWpTLmEb0.jpg",
	"ksSDzn7MEj.jpg",
	"mCzx91Sii1.jpg",
}

func init() {
	logger.Setup(true)
}

// createTestDirectory creates and fills a directory with some example content.
func createTestDirectory(directory string) error {
	if err := os.MkdirAll(directory, 0755); err != nil {
		return err
	}

	return testutil.FillDirectory(directory, testutil.FileTree{
		Dirs: map[string]*testutil.FileTree{
			"nested": {
				Files: []string{"ignore.txt"},
			},
		},
		Files: fileFixtures,
	})
}

func TestGetUnprocessedImages(t *testing.T) {
	var dir = "testdata"
	var images []model.Image
	var err error
	var expected []model.Image
	var imageFixtures []model.Image

	configuration := config.Load()
	configuration.Images = dir
	unprocessedImages := filepath.Join(dir, configuration.UnprocessedImagesFolder)

	for _, v := range fileFixtures {
		imageFixtures = append(imageFixtures, model.Image{
			File:               filepath.Join(configuration.UnprocessedImagesFolder, v),
			ProposedCategories: []string{},
			AssignedCategories: []string{},
			StarredCategory:    nil,
		})
	}

	defer os.RemoveAll(dir)
	if err := createTestDirectory(unprocessedImages); err != nil {
		format, args := testutil.FormatTestError(
			"Unable to create file fixtures.",
			map[string]interface{}{
				"error": err,
			})
		t.Fatalf(format, args...)
	}

	expected = imageFixtures
	images, err = controller.GetUnprocessedImages(model.ImageOptions{})
	if !reflect.DeepEqual(images, expected) {
		format, args := testutil.FormatTestError(
			"Returned images do not match expectations.",
			map[string]interface{}{
				"expected": expected,
				"got":      images,
				"error":    err,
			})
		t.Errorf(format, args...)
	}

	expected = []model.Image{
		imageFixtures[0],
		imageFixtures[1],
		imageFixtures[2],
	}
	images, err = controller.GetUnprocessedImages(model.ImageOptions{Count: util.IntPtr(3)})
	if !reflect.DeepEqual(images, expected) {
		format, args := testutil.FormatTestError(
			"Returned images do not match expectations.",
			map[string]interface{}{
				"expected": expected,
				"got":      images,
				"error":    err,
			})
		t.Errorf(format, args...)
	}

	expected = []model.Image{
		imageFixtures[3],
		imageFixtures[4],
	}
	images, err = controller.GetUnprocessedImages(model.ImageOptions{Count: util.IntPtr(5), LastImage: &fileFixtures[2]})
	if !reflect.DeepEqual(images, expected) {
		format, args := testutil.FormatTestError(
			"Returned images do not match expectations.",
			map[string]interface{}{
				"expected": expected,
				"got":      images,
				"error":    err,
			})
		t.Errorf(format, args...)
	}
}

func TestUpsertImage(t *testing.T) {
	var dir = "testdata"

	configuration := config.Load()
	configuration.Images = dir
	unprocessedImages := filepath.Join(dir, configuration.UnprocessedImagesFolder)

	defer os.RemoveAll(dir)
	mongodb.Connect(context.Background(), fmt.Sprintf(`mongodb://%s`, configuration.DatabaseHost))
	defer testutil.CleanCollection(t, configuration.Database, "image")

	if err := os.MkdirAll(unprocessedImages, 0755); err != nil {
		t.Fatal("Unable to create the unprocessed image folder.", err)
	}

	if err := testutil.TouchFile(
		filepath.Join(unprocessedImages, "test.jpg"),
	); err != nil {
		t.Fatal("Unable to create the test image.", err)
	}

	expected := model.Image{
		File:               filepath.Join(configuration.ProcessedImagesFolder, "test.jpg"),
		ProposedCategories: []string{},
		AssignedCategories: []string{"Category 1"},
		StarredCategory:    util.StringPtr("Category 1"),
	}
	image, err := controller.UpsertImage(model.Image{
		File:               filepath.Join(configuration.UnprocessedImagesFolder, "test.jpg"),
		AssignedCategories: []string{},
		ProposedCategories: []string{},
		StarredCategory:    util.StringPtr("Category 1"),
	})

	if err != nil || !reflect.DeepEqual(*image, expected) {
		format, args := testutil.FormatTestError(
			"Inserted image does not match expectations.",
			map[string]interface{}{
				"expected": expected,
				"got":      image,
				"error":    err,
			})
		t.Errorf(format, args...)
	}
}
