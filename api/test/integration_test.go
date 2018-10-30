package test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"

	"tagallery.com/api/config"
	"tagallery.com/api/model"
	"tagallery.com/api/mongodb"
	"tagallery.com/api/router"
	"tagallery.com/api/testutil"
	"tagallery.com/api/util"
)

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
	return nil
}

// createFileFixtures creates and fills a directory with some example content.
func createFileFixtures() error {

	return testutil.FillDirectory(config.GetConfig().Unprocessed_Images, testutil.FileTree{
		Dirs: map[string]*testutil.FileTree{
			"nested": &testutil.FileTree{
				Files: []string{"ignore.jpg"},
			},
		},
		Files: fileFixtures,
	})
}

// createDBFixtures inserts a set of test images into the database.
func createDBFixtures() error {
	images := make([]interface{}, len(imageFixtures))
	for k, v := range imageFixtures {
		images[k] = v
	}

	return testutil.InsertIntoMongoDb(
		config.GetConfig().Database_Host,
		config.GetConfig().Database,
		"images", images,
	)
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
	// Manually set the config here for debugging without a config file or env vars
	// config.SetConfig(config.Configuration{
	// 	Database:           "tagallery",
	// 	Database_Host:      "localhost:27017",
	// 	Unprocessed_Images: "/home/catyphram/Projects/Tagallery/code/api/src/test/images/",
	// 	Port:               3333,
	// })
	config.Load()
}

func TestAPI(t *testing.T) {
	if err := createDBFixtures(); err != nil {
		format, args := testutil.FormatTestError(
			"An error while creating the database fixtures has occured.",
			map[string]interface{}{
				"error": err,
			})
		t.Errorf(format, args...)
		return
	}
	defer dropDb(t)
	if err := createFileFixtures(); err != nil {
		format, args := testutil.FormatTestError(
			"An error while creating the file fixtures has occured.",
			map[string]interface{}{
				"error": err,
			})
		t.Errorf(format, args...)
		return
	}
	defer util.EmptyDirectory(config.GetConfig().Unprocessed_Images)
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

	if !reflect.DeepEqual(categories, model.Categories) {
		format, args := testutil.FormatTestError(
			"Returned categories do not match expectations.",
			map[string]interface{}{
				"expected": model.Categories,
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
