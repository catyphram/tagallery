package controller

import (
	"os"
	"reflect"
	"testing"

	"tagallery.com/api/config"
	"tagallery.com/api/model"
	"tagallery.com/api/testutil"
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

// createTestDirectory creates and fills a directory with some example content.
func createTestDirectory(directory string) error {

	if err := os.Mkdir(directory, os.ModePerm); err != nil {
		return err
	}

	return testutil.FillDirectory(directory, testutil.FileTree{
		Dirs: map[string]*testutil.FileTree{
			"nested": &testutil.FileTree{
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

	config.SetConfig(config.Configuration{
		Unprocessed_Images: dir,
	})

	createTestDirectory(dir)

	expected = []model.Image{
		{File: fileFixtures[0]},
		{File: fileFixtures[1]},
		{File: fileFixtures[2]},
		{File: fileFixtures[3]},
		{File: fileFixtures[4]},
	}
	images, err = GetUnprocessedImages(5, "")
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
		{File: fileFixtures[0]},
		{File: fileFixtures[1]},
		{File: fileFixtures[2]},
	}
	images, err = GetUnprocessedImages(3, "")
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
		{File: fileFixtures[3]},
		{File: fileFixtures[4]},
	}
	images, err = GetUnprocessedImages(5, fileFixtures[2])
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

	os.RemoveAll(dir)
}
