package test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
	"time"

	"tagallery.com/api/config"
	"tagallery.com/api/model"
	"tagallery.com/api/mongodb"
	"tagallery.com/api/router"
	"tagallery.com/api/testutil"
	"tagallery.com/api/util"
)

var categoryFixtures = []model.Category{
	{Name: "Category 1", Key: "category1", Description: "Category 1 description."},
	{Name: "Category 2", Key: "category2", Description: "Category 2 description."},
	{Name: "Category 3", Key: "category3", Description: "Category 3 description."},
	{Name: "Category 4", Key: "category4", Description: "Category 4 description."},
	{Name: "Category 5", Key: "category5", Description: "Category 5 description."},
	{Name: "Category 6", Key: "category6", Description: "Category 6 description."},
}

var imageFixtures = []model.Image{
	{File: "kQqIEGzwy6.jpg"},
	{File: "QLRLhnOrzA.jpg", AssignedCategories: []string{"Category 1", "Category 2"}},
	{File: "P0JDfeSr1v.jpg", ProposedCategories: []string{"Category 2"}},
	{File: "gH3sY4B2Wj.png", StarredCategory: "Category 1"},
	{File: "kl7RdVTAhp.jpg"},
	{File: "SRymDEkNbW.jpg", AssignedCategories: []string{"Category 2"}, ProposedCategories: []string{"Category 3"}},
	{File: "HnIPvWuPsI.jpg", ProposedCategories: []string{"Category 1", "Category 3"}},
	{File: "f.jpg"},
	{File: "obDrDA17Q9.jpg"},
	{File: "DiKkLH5kmf.jpg", AssignedCategories: []string{"Category 3"}},
	{File: "Je1p4msN6N.jpg", ProposedCategories: []string{"Category 2"}},
	{File: "8fNNuCssYi.png", StarredCategory: "Category 1"},
	{File: "sHJODFlvHI.jpg", AssignedCategories: []string{"Category 1"}},
	{File: "MF3onraKKV.png", AssignedCategories: []string{"Category 2"}, ProposedCategories: []string{"Category 1", "Category 3"}},
	{File: "O2hCNHYgYp.jpg", ProposedCategories: []string{"Category 1"}, StarredCategory: "Category 2"},
	{File: "N3Ynmj4GWt.jpg"},
	{File: "k0sCGdUC1F.jpg", StarredCategory: "Category 2"},
	{File: "GB3kVzAvfq", AssignedCategories: []string{"Category 2"}, ProposedCategories: []string{"Category 1", "Category 2"}},
	{File: "cUanJL8LAs.jpg", StarredCategory: "Category 2"},
	{File: "4jZUXo1cuG.jpg"},
}

// fileFixtures has to be sorted alphabetically to match the order of the returned images of the API.
var fileFixtures = []string{
	"0pFeM5WJIw.jpg",
	"12mlTsfTGX.jpg",
	"2AJm5tDwdh.jpg",
	"5JXQjezrHT.png",
	"6gfvyeYRXN.jpg",
	"6ld9fwduld.jpg",
	"A2SLbLVO1r.jpg",
	"EnCMYLR2LQ.jpg",
	"K8A4MYH0fB.jpg",
	"KlFt0RXsJC.jpg",
	"NYfxTaWQFx.jpg",
	"PzxY6JYRQp.png",
	"ReD4LDGQKk.png",
	"YIWj1hdgP6.jpg",
	"gPptk1Svj4.jpg",
	"l9C4bba1ys.jpg",
	"mM99bOCm68.jpg",
	"mO8q7Q87Hs.jpg",
	"rYqBjV4BvA.jpg",
	"rmQin0A17h.jpg",
}

// startAPI boots to API.
func startAPI() error {
	_, _, err := mongodb.Init(mongodb.DatabaseOptions{
		Database: config.GetConfig().Database,
		Host:     config.GetConfig().Database_Host,
	})
	if err != nil {
		return err
	}

	go http.ListenAndServe(fmt.Sprintf(":%v", config.GetConfig().Port), router.CreateRouter())

	// Wait 100 milliseconds to ensure that the http server has started
	// before the integration tests are run. A race condition may otherwise occur.
	time.Sleep(100 * time.Millisecond)

	return nil
}

// createFileFixtures creates and fills a directory with some example content.
func createFileFixtures(directory string) error {

	return testutil.FillDirectory(directory, testutil.FileTree{
		Dirs: map[string]*testutil.FileTree{
			"nested": {
				Files: []string{"ignore.jpg"},
			},
		},
		Files: fileFixtures,
	})
}

// createDBFixtures inserts a set of test images into the database.
func createDBFixtures() error {
    host := config.GetConfig().Database_Host
    database := config.GetConfig().Database

    categories := make([]interface{}, len(categoryFixtures))
    for k, v := range categoryFixtures {
        categories[k] = v
    }

	images := make([]interface{}, len(imageFixtures))
	for k, v := range imageFixtures {
		images[k] = v
    }

	if err := testutil.InsertIntoMongoDb(
		host, database,
		"categories", categories,
	); err != nil {
		return err
	}
    
	if err := testutil.InsertIntoMongoDb(
		host, database,
		"images", images,
	); err != nil {
		return err
	}

	return nil
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

func init() {
	config.Load()

	// Manually set the config here for debugging without a config file or env vars
	// config.SetConfig(config.Configuration{
	// 	Database:           "tagallery",
	// 	Database_Host:      "localhost:27017",
	//  Unprocessed_Images: "testdata",
	// 	Port:               3333,
	// })

}

func TestAPI(t *testing.T) {

    var directory = config.GetConfig().Unprocessed_Images
    
    defer dropDb(t)    

	if err := createDBFixtures(); err != nil {
		format, args := testutil.FormatTestError(
			"An error while creating the database fixtures has occured.",
			map[string]interface{}{
				"error": err,
			})
		t.Errorf(format, args...)
		return
	}

	if err := createFileFixtures(directory); err != nil {
		format, args := testutil.FormatTestError(
			"An error while creating the file fixtures has occured.",
			map[string]interface{}{
				"error": err,
			})
		t.Errorf(format, args...)
		return
	}
	defer util.EmptyDirectory(directory)

	if err := startAPI(); err != nil {
		format, args := testutil.FormatTestError(
			"An error while starting the API has occured.",
			map[string]interface{}{
				"error": err,
			})
		t.Errorf(format, args...)
		return
	}

	t.Run("GetCategories", testGetCategories)
	t.Run("GetImages", testGetImages)

}

// getRequest sends a HTTP request to {route} and parses the returned data into the type of {response}.
// The server host and port will be prefixed automatically.
func getRequest(route string, response interface{}) error {

	url := fmt.Sprintf("http://localhost:%v%v", config.GetConfig().Port, route)

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	json.Unmarshal(body, response)
	return nil
}

func testGetCategories(t *testing.T) {

	var categories []model.Category

	if err := getRequest("/categories", &categories); err != nil {
		format, args := testutil.FormatTestError(
			"Request failed.",
			map[string]interface{}{
				"error": err,
			})
		t.Errorf(format, args...)
	}

	if !reflect.DeepEqual(categories, categoryFixtures) {
		format, args := testutil.FormatTestError(
			"Returned categories do not match expectations.",
			map[string]interface{}{
				"expected": categoryFixtures,
				"got":      categories,
			})
		t.Errorf(format, args...)
	}
}

func testGetImages(t *testing.T) {
	var images []model.Image
	var expected []model.Image

	for i := 0; i < 15; i++ {
		expected = append(expected, imageFixtures[i])
	}
	if err := getRequest("/images", &images); err != nil {
		format, args := testutil.FormatTestError(
			"Request failed.",
			map[string]interface{}{
				"error": err,
			})
		t.Errorf(format, args...)
	}
	if !reflect.DeepEqual(expected, images) {
		format, args := testutil.FormatTestError(
			"Returned images do not match expectations.",
			map[string]interface{}{
				"expected": expected,
				"got":      images,
			})
		t.Errorf(format, args...)
	}

	expected = []model.Image{}
	for i := 0; i < 10; i++ {
		expected = append(expected, imageFixtures[i])
	}
	if err := getRequest("/images?count=10", &images); err != nil {
		format, args := testutil.FormatTestError(
			"Request failed.",
			map[string]interface{}{
				"error": err,
			})
		t.Errorf(format, args...)
	}
	if !reflect.DeepEqual(expected, images) {
		format, args := testutil.FormatTestError(
			"Returned images do not match expectations.",
			map[string]interface{}{
				"expected": expected,
				"got":      images,
			})
		t.Errorf(format, args...)
	}

	expected = []model.Image{
		imageFixtures[1],
	}
	if err := getRequest("/images?categories=Category+1&categories=Category+2", &images); err != nil {
		format, args := testutil.FormatTestError(
			"Request failed.",
			map[string]interface{}{
				"error": err,
			})
		t.Errorf(format, args...)
	}
	if !reflect.DeepEqual(expected, images) {
		format, args := testutil.FormatTestError(
			"Returned images do not match expectations.",
			map[string]interface{}{
				"expected": expected,
				"got":      images,
			})
		t.Errorf(format, args...)
	}

	expected = []model.Image{
		imageFixtures[12],
	}
	if err := getRequest(fmt.Sprintf("/images?categories=Category+1&last=%v", imageFixtures[1].File), &images); err != nil {
		format, args := testutil.FormatTestError(
			"Request failed.",
			map[string]interface{}{
				"error": err,
			})
		t.Errorf(format, args...)
	}
	if !reflect.DeepEqual(expected, images) {
		format, args := testutil.FormatTestError(
			"Returned images do not match expectations.",
			map[string]interface{}{
				"expected": expected,
				"got":      images,
			})
		t.Errorf(format, args...)
	}

	expected = []model.Image{
		imageFixtures[12],
	}
	if err := getRequest(fmt.Sprintf("/images?categories=Category+1&last=%v", imageFixtures[1].File), &images); err != nil {
		format, args := testutil.FormatTestError(
			"Request failed.",
			map[string]interface{}{
				"error": err,
			})
		t.Errorf(format, args...)
	}
	if !reflect.DeepEqual(expected, images) {
		format, args := testutil.FormatTestError(
			"Returned images do not match expectations.",
			map[string]interface{}{
				"expected": expected,
				"got":      images,
			})
		t.Errorf(format, args...)
	}

	expected = []model.Image{
		imageFixtures[2],
		imageFixtures[10],
		imageFixtures[17],
	}
	if err := getRequest("/images?status=autocategorized&categories=Category+2", &images); err != nil {
		format, args := testutil.FormatTestError(
			"Request failed.",
			map[string]interface{}{
				"error": err,
			})
		t.Errorf(format, args...)
	}
	if !reflect.DeepEqual(expected, images) {
		format, args := testutil.FormatTestError(
			"Returned images do not match expectations.",
			map[string]interface{}{
				"expected": expected,
				"got":      images,
			})
		t.Errorf(format, args...)
	}

	expected = []model.Image{
		imageFixtures[0],
		imageFixtures[4],
		imageFixtures[7],
		imageFixtures[8],
		imageFixtures[15],
		imageFixtures[19],
	}
	if err := getRequest("/images?status=uncategorized", &images); err != nil {
		format, args := testutil.FormatTestError(
			"Request failed.",
			map[string]interface{}{
				"error": err,
			})
		t.Errorf(format, args...)
	}
	if !reflect.DeepEqual(expected, images) {
		format, args := testutil.FormatTestError(
			"Returned images do not match expectations.",
			map[string]interface{}{
				"expected": expected,
				"got":      images,
			})
		t.Errorf(format, args...)
	}

	expected = []model.Image{}
	for i := 0; i < 15; i++ {
		expected = append(expected, model.Image{
			File: fileFixtures[i],
		})
	}
	if err := getRequest("/images?status=unprocessed", &images); err != nil {
		format, args := testutil.FormatTestError(
			"Request failed.",
			map[string]interface{}{
				"error": err,
			})
		t.Errorf(format, args...)
	}
	if !reflect.DeepEqual(expected, images) {
		format, args := testutil.FormatTestError(
			"Returned images do not match expectations.",
			map[string]interface{}{
				"expected": expected,
				"got":      images,
			})
		t.Errorf(format, args...)
	}

	expected = []model.Image{}
	for i := 7; i < 12; i++ {
		expected = append(expected, model.Image{
			File: fileFixtures[i],
		})
	}
	if err := getRequest(fmt.Sprintf("/images?status=unprocessed&last=%v&count=5", fileFixtures[6]), &images); err != nil {
		format, args := testutil.FormatTestError(
			"Request failed.",
			map[string]interface{}{
				"error": err,
			})
		t.Errorf(format, args...)
	}
	if !reflect.DeepEqual(expected, images) {
		format, args := testutil.FormatTestError(
			"Returned images do not match expectations.",
			map[string]interface{}{
				"expected": expected,
				"got":      images,
			})
		t.Errorf(format, args...)
	}

	expected = []model.Image{}
	for i := 7; i < 12; i++ {
		expected = append(expected, model.Image{
			File: fileFixtures[i],
		})
	}
	if err := getRequest(fmt.Sprintf("/images?status=unprocessed&last=%v&count=5", fileFixtures[6]), &images); err != nil {
		format, args := testutil.FormatTestError(
			"Request failed.",
			map[string]interface{}{
				"error": err,
			})
		t.Errorf(format, args...)
	}
	if !reflect.DeepEqual(expected, images) {
		format, args := testutil.FormatTestError(
			"Returned images do not match expectations.",
			map[string]interface{}{
				"expected": expected,
				"got":      images,
			})
		t.Errorf(format, args...)
	}
}
