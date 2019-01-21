package mongodb

import (
	"reflect"
	"testing"

	"tagallery.com/api/config"
	"tagallery.com/api/model"
	"tagallery.com/api/testutil"
)

var categoryFixtures = []model.Category{
	{Name: "Test 1", Key: "test1", Description: "Test 1 Category"},
	{Name: "Test 2", Key: "test2", Description: "Test 2 Category"},
	{Name: "Test 3", Key: "test3", Description: "Test 3 Category"},
	{Name: "Test 4", Key: "test4", Description: "Test 4 Category"},
	{Name: "Test 5", Key: "test5", Description: "Test 5 Category"},
	{Name: "Test 6", Key: "test6", Description: "Test 6 Category"},
}

// createCategoryFixtures creates a set of category entries in the database for testing.
func createCategoryFixtures(host string, database string) error {
	categories := make([]interface{}, len(categoryFixtures))
	for k, v := range categoryFixtures {
		categories[k] = v
	}

	return testutil.InsertIntoMongoDb(host, database, "categories", categories)
}

func init() {
	config.Load()

	// Manually set the config here for debugging without a config file or env vars
	config.SetConfig(config.Configuration{
		Database:      "tagallery",
		Database_Host: "localhost:27017",
	})
}

func TestGetCategories(t *testing.T) {
	var dbCategories []model.Category
	var expectedCategories []model.Category

	if err := createCategoryFixtures(
		config.GetConfig().Database_Host,
		config.GetConfig().Database,
	); err != nil {
		format, args := testutil.FormatTestError(
			"Unable to create category fixtures in the database.",
			map[string]interface{}{
				"error": err,
			})
		t.Errorf(format, args...)
		return
	}

	defer testutil.DropMongoDb(
		config.GetConfig().Database_Host,
		config.GetConfig().Database,
		t,
	)

	Init(DatabaseOptions{
		Host:     config.GetConfig().Database_Host,
		Database: config.GetConfig().Database,
	})

	expectedCategories = categoryFixtures
	dbCategories, _ = GetCategories()
	if !reflect.DeepEqual(dbCategories, expectedCategories) {
		format, args := testutil.FormatTestError(
			"Expected categorys from database to equal the fixture.", map[string]interface{}{
				"dbCategories":       dbCategories,
				"expectedCategories": expectedCategories,
			})
		t.Errorf(format, args...)
	}
}
