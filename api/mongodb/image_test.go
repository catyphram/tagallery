package mongodb

import (
	"reflect"
	"testing"

	"tagallery.com/api/config"
	"tagallery.com/api/model"
	"tagallery.com/api/testutil"
)

var imageFixtures = []model.Image{
	{File: "test1.jpg"},
	{File: "test2.jpg", AssignedCategories: []string{"Category 1", "Category 2"}},
	{File: "test3.jpg", ProposedCategories: []string{"Category 2"}},
	{File: "test4.jpg", StarredCategory: "Category 1"},
	{File: "test5.jpg", AssignedCategories: []string{"Category 2"}, ProposedCategories: []string{"Category 1", "Category 3"}},
	{File: "test6.jpg", ProposedCategories: []string{"Category 1"}, StarredCategory: "Category 2"},
}

// dropDb connects to and drops the database as specified in the config.
func dropDb(t *testing.T) {
	err := testutil.CleanMongoDb(
		config.GetConfig().Database_Host,
		config.GetConfig().Database,
	)
	if err != nil {
		format, args := testutil.FormatTestError(
			"Failed to clean up the database after testing.",
			map[string]interface{}{
				"error": err,
			})
		t.Errorf(format, args...)
	}
}

// createImageFixtures creates a set of image entries in the database for testing.
func createImageFixtures(host string, database string) error {
	images := make([]interface{}, len(imageFixtures))
	for k, v := range imageFixtures {
		images[k] = v
	}

	return testutil.InsertIntoMongoDb(host, database, "images", images)
}

func init() {
	// Manually set the config here for debugging without a config file or env vars
	// config.SetConfig(config.Configuration{
	// 	Database:      "tagallery",
	// 	Database_Host: "localhost:27017",
	// })
	config.Load()
}

func TestGetImages(t *testing.T) {
	var dbImages []model.Image
	var expectedImages []model.Image

	if err := createImageFixtures(
		config.GetConfig().Database_Host,
		config.GetConfig().Database,
	); err != nil {
		format, args := testutil.FormatTestError(
			"Unable to create image fixtures in the database.",
			map[string]interface{}{
				"error": err,
			})
		t.Errorf(format, args...)
		return
	}

	defer dropDb(t)

	Init(DatabaseOptions{
		Host:     config.GetConfig().Database_Host,
		Database: config.GetConfig().Database,
	})

	expectedImages = imageFixtures
	dbImages, _ = GetImages(10, &model.CategoryMap{}, "")
	if !reflect.DeepEqual(dbImages, expectedImages) {
		format, args := testutil.FormatTestError(
			"Expected images from database to equal the fixture.", map[string]interface{}{
				"dbImages":       dbImages,
				"expectedImages": expectedImages,
			})
		t.Errorf(format, args...)
	}

	expectedImages = []model.Image{
		imageFixtures[0],
		imageFixtures[1],
		imageFixtures[2],
	}
	dbImages, _ = GetImages(3, &model.CategoryMap{}, "")
	if !reflect.DeepEqual(dbImages, expectedImages) {
		format, args := testutil.FormatTestError(
			"Expected images from database to equal the fixture.", map[string]interface{}{
				"dbImages":       dbImages,
				"expectedImages": expectedImages,
			})
		t.Errorf(format, args...)
	}

	expectedImages = []model.Image{
		imageFixtures[3],
		imageFixtures[4],
		imageFixtures[5],
	}
	dbImages, _ = GetImages(10, &model.CategoryMap{}, imageFixtures[2].File)
	if !reflect.DeepEqual(dbImages, expectedImages) {
		format, args := testutil.FormatTestError(
			"Expected images from database to equal the fixture.", map[string]interface{}{
				"dbImages":       dbImages,
				"expectedImages": expectedImages,
			})
		t.Errorf(format, args...)
	}

	expectedImages = []model.Image{imageFixtures[0]}
	dbImages, _ = GetImages(10, nil, "")
	if !reflect.DeepEqual(dbImages, expectedImages) {
		format, args := testutil.FormatTestError(
			"Expected images from database to equal the fixture.", map[string]interface{}{
				"dbImages":       dbImages,
				"expectedImages": expectedImages,
			})
		t.Errorf(format, args...)
	}

	expectedImages = []model.Image{
		imageFixtures[1],
		imageFixtures[4],
	}
	dbImages, _ = GetImages(10, &model.CategoryMap{
		Assigned: []string{"Category 2"},
	}, "")
	if !reflect.DeepEqual(dbImages, expectedImages) {
		format, args := testutil.FormatTestError(
			"Expected images from database to equal the fixture.", map[string]interface{}{
				"dbImages":       dbImages,
				"expectedImages": expectedImages,
			})
		t.Errorf(format, args...)
	}

	expectedImages = []model.Image{
		imageFixtures[4],
		imageFixtures[5],
	}
	dbImages, _ = GetImages(10, &model.CategoryMap{
		Proposed: []string{"Category 1"},
	}, "")
	if !reflect.DeepEqual(dbImages, expectedImages) {
		format, args := testutil.FormatTestError(
			"Expected images from database to equal the fixture.", map[string]interface{}{
				"dbImages":       dbImages,
				"expectedImages": expectedImages,
			})
		t.Errorf(format, args...)
	}

	expectedImages = []model.Image{imageFixtures[3]}
	dbImages, _ = GetImages(10, &model.CategoryMap{
		Starred: "Category 1",
	}, "")
	if !reflect.DeepEqual(dbImages, expectedImages) {
		format, args := testutil.FormatTestError(
			"Expected images from database to equal the fixture.", map[string]interface{}{
				"dbImages":       dbImages,
				"expectedImages": expectedImages,
			})
		t.Errorf(format, args...)
	}
}

func TestUpsertImage(t *testing.T) {
	defer dropDb(t)

	Init(DatabaseOptions{
		Host:     config.GetConfig().Database_Host,
		Database: config.GetConfig().Database,
	})

	image := model.Image{
		File: "test",
	}
	err := UpsertImage(image)
	if err != nil {
		format, args := testutil.FormatTestError(
			"Expected image to be inserted, but got error.",
			map[string]interface{}{
				"error": err,
			})
		t.Errorf(format, args...)
	}
}
