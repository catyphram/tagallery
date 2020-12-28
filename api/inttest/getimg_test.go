package inttest

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"go.mongodb.org/mongo-driver/mongo/options"
	"tagallery.com/api/config"
	"tagallery.com/api/model"
	"tagallery.com/api/mongodb"
	"tagallery.com/api/testutil"
	"tagallery.com/api/util"
)

var processedImageFixtures = []model.Image{
	{File: "kQqIEGzwy6.jpg"},
	{File: "QLRLhnOrzA.jpg", AssignedCategories: []string{"Category 1", "Category 2"}},
	{File: "P0JDfeSr1v.jpg", ProposedCategories: []string{"Category 2"}},
	{File: "gH3sY4B2Wj.png", StarredCategory: util.StringPtr("Category 1")},
	{File: "kl7RdVTAhp.jpg"},
	{File: "SRymDEkNbW.jpg", AssignedCategories: []string{"Category 2"}, ProposedCategories: []string{"Category 3"}},
	{File: "HnIPvWuPsI.jpg", ProposedCategories: []string{"Category 1", "Category 3"}},
	{File: "f.jpg"},
	{File: "obDrDA17Q9.jpg"},
	{File: "DiKkLH5kmf.jpg", AssignedCategories: []string{"Category 3"}},
	{File: "Je1p4msN6N.jpg", ProposedCategories: []string{"Category 2"}},
	{File: "8fNNuCssYi.png", StarredCategory: util.StringPtr("Category 1")},
	{File: "sHJODFlvHI.jpg", AssignedCategories: []string{"Category 1"}},
	{File: "MF3onraKKV.png", AssignedCategories: []string{"Category 2"}, ProposedCategories: []string{"Category 1", "Category 3"}},
	{File: "O2hCNHYgYp.jpg", ProposedCategories: []string{"Category 1"}, StarredCategory: util.StringPtr("Category 2")},
	{File: "N3Ynmj4GWt.jpg"},
	{File: "k0sCGdUC1F.jpg", StarredCategory: util.StringPtr("Category 2")},
	{File: "GB3kVzAvfq", AssignedCategories: []string{"Category 2"}, ProposedCategories: []string{"Category 1", "Category 2"}},
	{File: "cUanJL8LAs.jpg", StarredCategory: util.StringPtr("Category 2")},
	{File: "4jZUXo1cuG.jpg"},
}

// unprocessedImageFixtures has to be sorted alphabetically to match the order of the returned images of the API.
var unprocessedImageFixtures = []string{
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

var unprocessedImages []string

// createFileFixtures creates and fills a directory with some example content.
func createFileFixtures(directory string, unprocessed []string, processed []string) error {
	var processedImages []string

	for _, v := range processedImageFixtures {
		processedImages = append(processedImages, v.File)
	}
	return testutil.FillDirectory(directory, testutil.FileTree{
		Dirs: map[string]*testutil.FileTree{
			config.Get().ProcessedImagesFolder: {
				Dirs: map[string]*testutil.FileTree{
					"nested": {
						Files: []string{"ignore.jpg"},
					},
				},
				Files: processed,
			},
			config.Get().UnprocessedImagesFolder: {
				Dirs: map[string]*testutil.FileTree{
					"nested": {
						Files: []string{"ignore.jpg"},
					},
				},
				Files: unprocessed,
			},
		},
	})
}

func createImageDBFixtures(ctx context.Context, db string) error {
	collection := mongodb.Client().Database(db).Collection("image")
	images := make([]interface{}, len(processedImageFixtures))
	for k, v := range processedImageFixtures {
		images[k] = v
	}

	if _, err := collection.InsertMany(ctx, images, options.InsertMany()); err != nil {
		return err
	}
	return nil
}

func getImagesSetup(t *testing.T) func() {
	var processedImages []string
	for _, v := range processedImageFixtures {
		processedImages = append(processedImages, v.File)
	}

	if err := createImageDBFixtures(context.Background(), config.Get().Database); err != nil {
		format, args := testutil.FormatTestError(
			"An error while creating the database fixtures has occured.",
			map[string]interface{}{
				"error": err,
			})
		t.Fatalf(format, args...)
	}

	_ = os.MkdirAll(config.Get().Images, 0755)
	if err := createFileFixtures(
		config.Get().Images,
		unprocessedImageFixtures,
		processedImages,
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

func GetImages(t *testing.T) {
	var images []model.Image
	var expected []model.Image

	for _, v := range unprocessedImageFixtures {
		unprocessedImages = append(unprocessedImages, filepath.Join(config.Get().UnprocessedImagesFolder, v))
	}

	defer getImagesSetup(t)()

	for i := 0; i < 15; i++ {
		expected = append(expected, processedImageFixtures[i])
	}
	if err := GetRequest(apiURL("/image"), &images); err != nil {
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
		expected = append(expected, processedImageFixtures[i])
	}
	if err := GetRequest(apiURL("/image?count=10"), &images); err != nil {
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
		processedImageFixtures[1],
	}
	if err := GetRequest(
		apiURL("/image?categories=Category+1&categories=Category+2"), &images,
	); err != nil {
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
		processedImageFixtures[12],
	}
	if err := GetRequest(apiURL(fmt.Sprintf(
		"/image?categories=Category+1&lastImage=%v",
		processedImageFixtures[1].File)), &images,
	); err != nil {
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
		processedImageFixtures[12],
	}
	if err := GetRequest(apiURL(fmt.Sprintf(
		"/image?categories=Category+1&lastImage=%v",
		processedImageFixtures[1].File)), &images,
	); err != nil {
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
		processedImageFixtures[2],
		processedImageFixtures[10],
		processedImageFixtures[17],
	}
	if err := GetRequest(apiURL("/image?status=autocategorized&categories=Category+2"), &images); err != nil {
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
		processedImageFixtures[0],
		processedImageFixtures[4],
		processedImageFixtures[7],
		processedImageFixtures[8],
		processedImageFixtures[15],
		processedImageFixtures[19],
	}
	if err := GetRequest(apiURL("/image?status=uncategorized"), &images); err != nil {
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
			File:               unprocessedImages[i],
			AssignedCategories: []string{},
			ProposedCategories: []string{},
		})
	}
	if err := GetRequest(apiURL("/image?status=unprocessed"), &images); err != nil {
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
			File:               unprocessedImages[i],
			AssignedCategories: []string{},
			ProposedCategories: []string{},
		})
	}
	if err := GetRequest(apiURL(fmt.Sprintf(
		"/image?status=unprocessed&lastImage=%v&count=5",
		unprocessedImageFixtures[6])), &images,
	); err != nil {
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
			File:               unprocessedImages[i],
			AssignedCategories: []string{},
			ProposedCategories: []string{},
		})
	}
	if err := GetRequest(apiURL(fmt.Sprintf(
		"/image?status=unprocessed&lastImage=%v&count=5",
		unprocessedImageFixtures[6])), &images,
	); err != nil {
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
