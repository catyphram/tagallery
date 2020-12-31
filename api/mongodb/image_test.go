package mongodb_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"go.mongodb.org/mongo-driver/mongo/options"
	"tagallery.com/api/config"
	"tagallery.com/api/logger"
	"tagallery.com/api/model"
	"tagallery.com/api/mongodb"
	"tagallery.com/api/testutil"
	"tagallery.com/api/util"
)

func init() {
	logger.Setup(true)
}

var imageFixtures = []model.Image{
	{File: "test1.jpg"},
	{File: "test2.jpg", AssignedCategories: []string{"Category 1", "Category 2"}},
	{File: "test3.jpg", ProposedCategories: []string{"Category 2"}},
	{File: "test4.jpg", StarredCategory: util.StringPtr("Category 1")},
	{File: "test5.jpg", AssignedCategories: []string{"Category 2"}, ProposedCategories: []string{"Category 1", "Category 3"}},
	{File: "test6.jpg", ProposedCategories: []string{"Category 1"}, StarredCategory: util.StringPtr("Category 2")},
}

// createImageFixtures inserts the image fixtures into the database.
func createImageFixtures(ctx context.Context, db string) error {
	collection := mongodb.Client().Database(db).Collection("image")

	images := make([]interface{}, len(imageFixtures))
	for k, v := range imageFixtures {
		images[k] = v
	}

	_, err := collection.InsertMany(ctx, images, options.InsertMany())

	return err
}

func TestGetImages(t *testing.T) {
	var dbImages []model.Image
	var expectedImages []model.Image

	configuration := config.Load()

	mongodb.Connect(context.Background(), fmt.Sprintf(`mongodb://%s`, configuration.DatabaseHost))
	defer testutil.CleanCollection(t, configuration.Database, "image")

	if err := createImageFixtures(context.Background(), configuration.Database); err != nil {
		format, args := testutil.FormatTestError(
			"Unable to create image fixtures in the database.",
			map[string]interface{}{
				"error": err,
			})
		t.Fatalf(format, args...)
	}

	expectedImages = imageFixtures
	dbImages, _ = mongodb.GetImages(
		model.ImageOptions{Count: util.IntPtr(10)},
		&model.CategoryMap{},
	)

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
	dbImages, _ = mongodb.GetImages(
		model.ImageOptions{Count: util.IntPtr(3)},
		&model.CategoryMap{},
	)
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
	dbImages, _ = mongodb.GetImages(
		model.ImageOptions{
			Count:     util.IntPtr(10),
			LastImage: util.StringPtr(imageFixtures[2].File),
		},
		&model.CategoryMap{},
	)
	if !reflect.DeepEqual(dbImages, expectedImages) {
		format, args := testutil.FormatTestError(
			"Expected images from database to equal the fixture.", map[string]interface{}{
				"dbImages":       dbImages,
				"expectedImages": expectedImages,
			})
		t.Errorf(format, args...)
	}

	expectedImages = []model.Image{imageFixtures[0]}
	dbImages, _ = mongodb.GetImages(
		model.ImageOptions{Count: util.IntPtr(10)},
		nil,
	)
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
	dbImages, _ = mongodb.GetImages(
		model.ImageOptions{Count: util.IntPtr(10)},
		&model.CategoryMap{
			Assigned: []string{"Category 2"},
		},
	)
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
	dbImages, _ = mongodb.GetImages(
		model.ImageOptions{Count: util.IntPtr(10)},
		&model.CategoryMap{
			Assigned: []string{},
		},
	)
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
	dbImages, _ = mongodb.GetImages(
		model.ImageOptions{Count: util.IntPtr(10)},
		&model.CategoryMap{
			Proposed: []string{"Category 1"},
		},
	)
	if !reflect.DeepEqual(dbImages, expectedImages) {
		format, args := testutil.FormatTestError(
			"Expected images from database to equal the fixture.", map[string]interface{}{
				"dbImages":       dbImages,
				"expectedImages": expectedImages,
			})
		t.Errorf(format, args...)
	}

	expectedImages = []model.Image{
		imageFixtures[2],
		imageFixtures[4],
		imageFixtures[5],
	}
	dbImages, _ = mongodb.GetImages(
		model.ImageOptions{Count: util.IntPtr(10)},
		&model.CategoryMap{
			Proposed: []string{},
		},
	)
	if !reflect.DeepEqual(dbImages, expectedImages) {
		format, args := testutil.FormatTestError(
			"Expected images from database to equal the fixture.", map[string]interface{}{
				"dbImages":       dbImages,
				"expectedImages": expectedImages,
			})
		t.Errorf(format, args...)
	}

	expectedImages = []model.Image{imageFixtures[3]}
	dbImages, _ = mongodb.GetImages(
		model.ImageOptions{Count: util.IntPtr(10)},
		&model.CategoryMap{
			Starred: util.StringPtr("Category 1"),
		},
	)
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
	configuration := config.Load()

	image := model.Image{
		File: "test",
	}

	mongodb.Connect(context.Background(), fmt.Sprintf(`mongodb://%s`, configuration.DatabaseHost))
	defer testutil.CleanCollection(t, configuration.Database, "image")

	err := mongodb.UpsertImage(image)
	if err != nil {
		format, args := testutil.FormatTestError(
			"Expected image to be inserted.",
			map[string]interface{}{
				"error": err,
			})
		t.Errorf(format, args...)
	}
}
